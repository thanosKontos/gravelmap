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
	edgeFromFile  *os.File
	edgeToFile    *os.File
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
		edgeFromFile:  efF,
		edgeToFile:    etF,
		polylinesFile: plF,
	}, nil
}

func (fr *fileRead) Read(edges []gravelmap.Edge) ([]gravelmap.PresentableWay, error) {
	var presentableWays []gravelmap.PresentableWay

	for _, edge := range edges {
		nodeStart, err := fr.readEdgeStartFile(edge.NodeFrom)
		edgeToRec, err := fr.readEdgeToFile(*nodeStart, edge.NodeTo)
		if err != nil {
			return []gravelmap.PresentableWay{}, err
		}

		pl, err := fr.readPolylineFromFile(edgeToRec.polylinePosition)
		if err != nil {
			return []gravelmap.PresentableWay{}, err
		}

		presentableWays = append(presentableWays, gravelmap.PresentableWay{
			Polyline:    pl,
			SurfaceType: edgeToRec.wayType,
			ElevFrom:    edgeToRec.elevFrom,
			ElevTo:      edgeToRec.elevTo,
			Distance:    edgeToRec.distance,
			OsmID:       edgeToRec.osmID,
		})
	}

	fr.close()

	return presentableWays, nil
}

func (fr *fileRead) readEdgeToFile(edgeStart edgeStartRecord, edgeToId int32) (*edgeToRecord, error) {
	readOffset := edgeStart.NodeToOffset
	var distance, polylineLength int32
	var polylineOffset, osmID int64
	var wayType int8
	var elevationStart, elevationEnd int16
	found := false

	for i := 0; int8(i) < edgeStart.ConnectionsCnt; i++ {
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

			data = readNextBytes(fr.edgeToFile, 8)
			buffer = bytes.NewBuffer(data)
			binary.Read(buffer, binary.BigEndian, &osmID)

			found = true
			break
		}

		readOffset = readOffset + edgeToIndividualRecordSize
	}

	if !found {
		return nil, errors.New("edge to record not found")
	}

	edgeToRecord := edgeToRecord{
		nodeTo:           edgeToId,
		distance:         distance,
		wayType:          wayType,
		elevFrom:         elevationStart,
		elevTo:           elevationEnd,
		polylinePosition: polylinePosition{length: polylineLength, offset: polylineOffset},
		osmID:            osmID,
	}

	return &edgeToRecord, nil
}

func (fr *fileRead) readEdgeStartFile(edgeStartId int32) (*edgeStartRecord, error) {
	edgeStart := edgeStartRecord{}

	fr.edgeFromFile.Seek(int64(edgeStartId*edgeStartRecordSize), 0)
	data := readNextBytes(fr.edgeFromFile, edgeStartRecordSize)
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.BigEndian, &edgeStart); err != nil {
		return nil, err
	}

	return &edgeStart, nil
}

func (fr *fileRead) readPolylineFromFile(polylinePos polylinePosition) (string, error) {
	fr.polylinesFile.Seek(polylinePos.offset, 0)

	buffer := make([]byte, polylinePos.length)
	if _, err := fr.polylinesFile.Read(buffer); err != nil {
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
