package way

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"

	"github.com/thanosKontos/gravelmap"
)

type fileRead struct {
	storageDir string
}

func NewWayFileRead(storageDir string) *fileRead {
	return &fileRead{
		storageDir: storageDir,
	}
}

func (fr *fileRead) Read(ways []gravelmap.Way) ([]gravelmap.PresentableWay, error) {
	var presentableWays []gravelmap.PresentableWay

	esFl, err := os.Open(fmt.Sprintf("%s/%s", fr.storageDir, edgeStartFilename))
	defer esFl.Close()
	if err != nil {
		return []gravelmap.PresentableWay{}, err
	}

	plLkFl, err := os.Open(fmt.Sprintf("%s/%s", fr.storageDir, edgeToFilename))
	defer plLkFl.Close()
	if err != nil {
		return []gravelmap.PresentableWay{}, err
	}

	plFl, err := os.Open(fmt.Sprintf("%s/%s", fr.storageDir, polylinesFilename))
	defer plLkFl.Close()
	if err != nil {
		return []gravelmap.PresentableWay{}, err
	}

	//presentableWaysvar polylines []string
	for _, way := range ways {
		nodeStart, err := fr.readEdgeStartFile(esFl, way.EdgeFrom)

		edgeToRec, err := fr.readEdgeToFile(plLkFl, *nodeStart, way.EdgeTo)
		if err != nil {
			return []gravelmap.PresentableWay{}, err
		}

		pl, err := fr.readPolylineFromFile(plFl, edgeToRec.polylinePosition.length, edgeToRec.polylinePosition.offset)

		if err != nil {
			return []gravelmap.PresentableWay{}, err
		}

		presentableWays = append(presentableWays, gravelmap.PresentableWay{Polyline: pl, SurfaceType: edgeToRec.wayType, ElevationGrade: edgeToRec.grade})
	}

	return presentableWays, nil
}

func (fr *fileRead) readEdgeToFile(f *os.File, edgeStart edgeStartRecord, edgeToId int32) (*edgeToRecord, error) {
	readOffset := edgeStart.NodeToOffset
	var polylineLength int32
	var polylineOffset int64
	var wayType int8
	var grade float32
	found := false

	for i := 0; int32(i) < edgeStart.ConnectionsCnt; i++ {
		f.Seek(readOffset, 0)

		var storedEdgeTo int32

		data := readNextBytes(f, 4)
		buffer := bytes.NewBuffer(data)
		binary.Read(buffer, binary.BigEndian, &storedEdgeTo)

		if storedEdgeTo == edgeToId {
			data := readNextBytes(f, 1)
			buffer := bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &wayType)

			data = readNextBytes(f, 4)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &grade)

			data = readNextBytes(f, 4)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &polylineLength)

			data = readNextBytes(f, 8)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &polylineOffset)

			found = true
			break
		}

		readOffset = readOffset + edgeToIndividualRecordSize
	}

	if !found {
		return nil, errors.New("edge to record not found")
	}

	edgeToRecord := edgeToRecord{
		nodeTo: edgeToId,
		wayType: wayType,
		grade: grade,
		polylinePosition: polylinePosition{length: polylineLength, offset: polylineOffset},
	}

	return &edgeToRecord, nil
}

func (fr *fileRead) readPolylineFromFile(f *os.File, length int32, offset int64) (string, error) {
	f.Seek(offset, 0)

	buffer := make([]byte, length)
	_, err := f.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func (fr *fileRead) readEdgeStartFile(f *os.File, edgeStartId int32) (*edgeStartRecord, error) {
	edgeStart := edgeStartRecord{}

	f.Seek(int64(edgeStartId*edgeStartRecordSize), 0)

	data := readNextBytes(f, edgeStartRecordSize)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &edgeStart)
	if err != nil {
		return nil, err
	}

	return &edgeStart, nil
}

func readNextBytes(file *os.File, number int) []byte {
	byteSeq := make([]byte, number)
	_, _ = file.Read(byteSeq)

	return byteSeq
}


