package edge

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/thanosKontos/gravelmap"
)

const bBoxDir = "edge_bbox"

type bboxEdgeRecord struct {
	Pt gravelmap.Point
	EdgeID int32
}

type bboxFileStore struct {
	storageDir string
}

type bboxFileRead struct {
	storageDir string
	distanceCalc gravelmap.DistanceCalculator
}

func NewBBoxFileStore(storageDir string) *bboxFileStore {
	os.RemoveAll(fmt.Sprintf("%s/%s", storageDir, bBoxDir))
	os.Mkdir(fmt.Sprintf("%s/%s", storageDir, bBoxDir), 0777)

	return &bboxFileStore{
		storageDir: storageDir,
	}
}

func NewBBoxFileRead(storageDir string, dc gravelmap.DistanceCalculator) *bboxFileRead {
	return &bboxFileRead{
		storageDir: storageDir,
		distanceCalc: dc,
	}
}

func (fr *bboxFileRead) FindClosest(point gravelmap.Point) (int32, error) {
	filename := findBBoxFileFromPoint(point)

	f, err := os.Open(fmt.Sprintf("%s/%s/%s", fr.storageDir, bBoxDir, filename))
	defer f.Close()
	if err != nil {
		return 0, err
	}

	var closestEdge int32 = 0
	var closestEdgeDistance int32 = 0
	for {
		edgeRec := bboxEdgeRecord{}
		data := readNextBytes(f, 20)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &edgeRec)

		//fmt.Println(edgeRec)

		if closestEdge == 0 {
			closestEdge = edgeRec.EdgeID
			closestEdgeDistance = fr.distanceCalc.Calculate(edgeRec.Pt, point)
		} else {
			d := fr.distanceCalc.Calculate(edgeRec.Pt, point)
			if closestEdgeDistance > d {
				closestEdge = edgeRec.EdgeID
				closestEdgeDistance = d
			}
		}

		if edgeRec.EdgeID == 0 {
			if closestEdge == 0 {
				return 0, errors.New("no edge found")
			}

			break
		}
	}

	return closestEdge, nil
}

func readNextBytes(file *os.File, number int) []byte {
	byteSeq := make([]byte, number)
	_, _ = file.Read(byteSeq)

	return byteSeq
}

func (fs *bboxFileStore) BatchStore(ndBatch []gravelmap.GMNode) error {
	ndBatchFileMap := map[string][]gravelmap.GMNode{}
	for _, gmNd := range ndBatch {
		filename := findBBoxFileFromPoint(gmNd.Point)
		ndBatchFileMap[filename] = append(ndBatchFileMap[filename], gmNd)
	}

	for filename, batch := range ndBatchFileMap {
		err := fs.writeBatchToFile(filename, batch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *bboxFileStore) writeBatchToFile(filename string, ndBatch []gravelmap.GMNode) error {
	f, err := os.OpenFile(
		fmt.Sprintf("%s/%s/%s", fs.storageDir, bBoxDir, filename),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0777,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err != nil {
		return err
	}

	var recs []bboxEdgeRecord
	for _, gmNd := range ndBatch {
		recs = append(recs, bboxEdgeRecord{gmNd.Point, gmNd.ID})
	}

	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, recs)
	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())
	return err
}

func findBBoxFileFromPoint(p gravelmap.Point) string {
	// Level 1 bbox
	n := math.Floor(p.Lat)
	e := math.Floor(p.Lng)

	// Level 2 bbox
	l2 := math.Floor((p.Lat-n)*10)

	return fmt.Sprintf("N%.0fE%.0f_%.0f.bin", n, e, l2)
}

