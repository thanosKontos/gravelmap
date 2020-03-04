package edge

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/thanosKontos/gravelmap"
)

type bboxFileRead struct {
	storageDir   string
	distanceCalc gravelmap.DistanceCalculator
}

func NewBBoxFileRead(storageDir string, dc gravelmap.DistanceCalculator) *bboxFileRead {
	return &bboxFileRead{
		storageDir:   storageDir,
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
	var closestEdgeDistance int64 = 0
	for {
		edgeRec := bboxEdgeRecord{}
		data := readNextBytes(f, 20)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &edgeRec)

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
