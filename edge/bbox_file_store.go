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

func (fs *bboxFileStore) BatchStore(ndBatch []gravelmap.Node) error {
	ndBatchFileMap := map[string][]gravelmap.Node{}
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

func (fs *bboxFileStore) writeBatchToFile(filename string, ndBatch []gravelmap.Node) error {
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
		recs = append(recs, bboxEdgeRecord{gmNd.Point, int32(gmNd.Id)})
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

