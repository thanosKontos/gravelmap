package edge

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/thanosKontos/gravelmap"
)

type bboxFileStore struct {
	storageDir string
}

func NewBBoxFileStore(storageDir string) *bboxFileStore {
	os.RemoveAll(fmt.Sprintf("%s/%s", storageDir, bBoxDir))
	os.Mkdir(fmt.Sprintf("%s/%s", storageDir, bBoxDir), 0777)

	return &bboxFileStore{
		storageDir: storageDir,
	}
}

func (fs *bboxFileStore) BatchStore(ndPts []gravelmap.NodePoint) error {
	ndBatchFileMap := map[string][]gravelmap.NodePoint{}
	for _, gmNd := range ndPts {
		filename := findBBoxFileFromPoint(gmNd.Pt)
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

func (fs *bboxFileStore) writeBatchToFile(filename string, ndPts []gravelmap.NodePoint) error {
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

	var recs []gravelmap.NodePoint
	for _, ndPt := range ndPts {
		recs = append(recs, gravelmap.NodePoint{NodeID: ndPt.NodeID, Pt: ndPt.Pt})
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
	l2 := math.Floor((p.Lat - n) * 10)

	return fmt.Sprintf("N%.0fE%.0f_%.0f.bin", n, e, l2)
}
