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
	edgeFromFile *os.File
	edgeToFile *os.File
	polylinesFile *os.File
}

func NewWayFileRead(storageDir string) (*fileRead, error) {
	efF, err := os.Open(fmt.Sprintf("%s/%s", storageDir, edgeStartFilename))
	if err != nil {
		return nil, err
	}

	etF, err := os.Open(fmt.Sprintf("%s/%s", storageDir, edgeToFilename))
	if err != nil {
		return nil, err
	}

	plF, err := os.Open(fmt.Sprintf("%s/%s", storageDir, polylinesFilename))
	if err != nil {
		return nil, err
	}

	return &fileRead{
		edgeFromFile: efF,
		edgeToFile: etF,
		polylinesFile: plF,
	}, nil
}

func (fr *fileRead) Read(ways []gravelmap.Way) ([]gravelmap.PresentableWay, error) {
	var presentableWays []gravelmap.PresentableWay

	for _, way := range ways {
		nodeStart, err := fr.readEdgeStartFile(way.EdgeFrom)

		edgeToRec, err := fr.readEdgeToFile(*nodeStart, way.EdgeTo)
		if err != nil {
			return []gravelmap.PresentableWay{}, err
		}

		pl, err := fr.readPolylineFromFile(edgeToRec.polylinePosition)

		if err != nil {
			return []gravelmap.PresentableWay{}, err
		}

		presentableWays = append(presentableWays, gravelmap.PresentableWay{
			Polyline: pl,
			SurfaceType: edgeToRec.wayType,
			ElevationInfo: edgeToRec.ElevationInfo,
			Distance: edgeToRec.distance,
		})
	}

	fr.close()

	return presentableWays, nil
}

func (fr *fileRead) readEdgeToFile(edgeStart edgeStartRecord, edgeToId int32) (*edgeToRecord, error) {
	readOffset := edgeStart.NodeToOffset
	var distance, polylineLength int32
	var polylineOffset int64
	var wayType int8
	var grade float32
	var elevationStart, elevationEnd int16
	found := false

	for i := 0; int32(i) < edgeStart.ConnectionsCnt; i++ {
		fr.edgeToFile.Seek(readOffset, 0)

		var storedEdgeTo int32

		data := readNextBytes(fr.edgeToFile, 4)
		buffer := bytes.NewBuffer(data)
		binary.Read(buffer, binary.BigEndian, &storedEdgeTo)

		if storedEdgeTo == edgeToId {
			data := readNextBytes(fr.edgeToFile, 4)
			buffer := bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &distance)

			data = readNextBytes(fr.edgeToFile, 1)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &wayType)

			data = readNextBytes(fr.edgeToFile, 4)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &grade)

			data = readNextBytes(fr.edgeToFile, 2)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &elevationStart)

			data = readNextBytes(fr.edgeToFile, 2)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &elevationEnd)

			data = readNextBytes(fr.edgeToFile, 4)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &polylineLength)

			data = readNextBytes(fr.edgeToFile, 8)
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
		distance: distance,
		wayType: wayType,
		ElevationInfo: gravelmap.ElevationInfo{
			Grade: grade,
			From: elevationStart,
			To: elevationEnd,
		},
		polylinePosition: polylinePosition{length: polylineLength, offset: polylineOffset},
	}

	return &edgeToRecord, nil
}

func (fr *fileRead) readEdgeStartFile(edgeStartId int32) (*edgeStartRecord, error) {
	edgeStart := edgeStartRecord{}

	fr.edgeFromFile.Seek(int64(edgeStartId*edgeStartRecordSize), 0)
	data := readNextBytes(fr.edgeFromFile, edgeStartRecordSize)
	buffer := bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.BigEndian, &edgeStart)
	if err != nil {
		return nil, err
	}

	return &edgeStart, nil
}

func (fr *fileRead) readPolylineFromFile(polylinePos polylinePosition) (string, error) {
	fr.polylinesFile.Seek(polylinePos.offset, 0)

	buffer := make([]byte, polylinePos.length)
	_, err := fr.polylinesFile.Read(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer), nil
}

func readNextBytes(file *os.File, number int) []byte {
	byteSeq := make([]byte, number)
	_, _ = file.Read(byteSeq)

	return byteSeq
}

func (fr *fileRead) close() {
	fr.edgeFromFile.Close()
	fr.edgeToFile.Close()
	fr.polylinesFile.Close()
}
