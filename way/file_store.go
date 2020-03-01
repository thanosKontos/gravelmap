package way

import (
	"bytes"
	"encoding/binary"
	"os"
	"sort"

	"fmt"

	"github.com/thanosKontos/gravelmap"
)

type fileStore struct {
	storageDir string
	polylinesFile *os.File
	edgeFromFile *os.File
	edgeToFile *os.File
}

func NewFileStore(storageDir string) *fileStore {
	efF, _ := os.Create(fmt.Sprintf("%s/%s", storageDir, edgeStartFilename))
	etF, _ := os.Create(fmt.Sprintf("%s/%s", storageDir, edgeToFilename))
	plF, _ := os.Create(fmt.Sprintf("%s/%s", storageDir, polylinesFilename))

	return &fileStore{
		storageDir: storageDir,
		edgeFromFile: efF,
		edgeToFile: etF,
		polylinesFile: plF,
	}
}

type edgeToRecord struct {
	nodeTo int32
	wayType int8
	grade float32
	polylineLen int32
	polylineOffset int64
}

func (fs *fileStore) Store(ways map[int][]gravelmap.WayTo) error {
	var gmNodeIdsSorted []int
	for k := range ways {
		gmNodeIdsSorted = append(gmNodeIdsSorted, k)
	}
	sort.Ints(gmNodeIdsSorted)

	var offset int64 = 0
	var polylineOffset int64 = 0

	for _, gmNdID := range gmNodeIdsSorted {
		way := ways[gmNdID]

		var polylines []string
		var edgeToRecords []edgeToRecord

		for i := 0; i < len(way); i++ {
			polylineLen := int32(len(way[i].Polyline))

			edgeToRec := edgeToRecord{
				nodeTo: int32(way[i].NdTo),
				wayType: gravelmap.WayTypeUnaved,
				grade: 5.0,
				polylineLen: polylineLen,
				polylineOffset: polylineOffset,
			}
			edgeToRecords = append(edgeToRecords, edgeToRec)

			polylineOffset += int64(polylineLen)

			polylines = append(polylines, way[i].Polyline)
		}
		err := fs.writePolylinesFile(polylines)
		if err != nil {
			return err
		}

		err = fs.writeEdgeToFile(edgeToRecords)
		if err != nil {
			return err
		}

		edgeStart := edgeStartRecord{int32(len(way)), offset}

		err = fs.writeEdgeFromFile(gmNdID, edgeStart)
		if err != nil {
			return err
		}

		offset += int64(len(way))*edgeToIndividualRecordSize
	}

	fs.close()

	return nil
}

func (fs *fileStore) writeEdgeToFile(plsLookup []edgeToRecord) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, plsLookup)
	if err != nil {
		return err
	}

	_, err = fs.edgeToFile.Write(buf.Bytes())
	return err
}

func (fs *fileStore) writePolylinesFile(pls []string) error {
	for _, pl := range pls {
		_, err := fs.polylinesFile.WriteString(pl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *fileStore) writeEdgeFromFile(edgeStartId int, edgeStart edgeStartRecord) error {
	fs.edgeFromFile.Seek(int64(edgeStartId*edgeStartRecordSize), 0)

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, edgeStart)
	if err != nil {
		return err
	}

	_, err = fs.edgeFromFile.Write(buf.Bytes())

	return err
}

func (fs *fileStore) close() {
	fs.edgeFromFile.Close()
	fs.edgeToFile.Close()
	fs.polylinesFile.Close()
}
