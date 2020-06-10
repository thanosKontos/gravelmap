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

	var closestNode int32 = 0
	var closestNodeDistance int64 = 0
	for {
		nodePoint := gravelmap.NodePoint{}
		data := readNextBytes(f, 20)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &nodePoint)

		if closestNode == 0 {
			closestNode = nodePoint.NodeID
			closestNodeDistance = fr.distanceCalc.Calculate(nodePoint.Pt, point)
		} else {
			d := fr.distanceCalc.Calculate(nodePoint.Pt, point)
			if closestNodeDistance > d {
				closestNode = nodePoint.NodeID
				closestNodeDistance = d
			}
		}

		if nodePoint.NodeID == 0 {
			if closestNode == 0 {
				return 0, errors.New("no node found")
			}

			break
		}
	}

	return closestNode, nil
}

func readNextBytes(file *os.File, number int) []byte {
	byteSeq := make([]byte, number)
	_, _ = file.Read(byteSeq)

	return byteSeq
}
