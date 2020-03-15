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
	polylinesFile *os.File
	edgeFromFile  *os.File
	edgeToFile    *os.File
	pointEncoder  gravelmap.Encoder
}

func NewFileStore(storageDir string, pointEncoder gravelmap.Encoder) *fileStore {
	efF, _ := os.Create(fmt.Sprintf("%s/%s", storageDir, edgeStartFilename))
	etF, _ := os.Create(fmt.Sprintf("%s/%s", storageDir, edgeToFilename))
	plF, _ := os.Create(fmt.Sprintf("%s/%s", storageDir, polylinesFilename))

	return &fileStore{
		edgeFromFile:  efF,
		edgeToFile:    etF,
		polylinesFile: plF,
		pointEncoder:  pointEncoder,
	}
}

func (fs *fileStore) Store(ways map[int]map[int]gravelmap.EvaluatedWay) error {
	var gmNodeIdsSorted []int
	for k := range ways {
		gmNodeIdsSorted = append(gmNodeIdsSorted, k)
	}
	sort.Ints(gmNodeIdsSorted)

	var offset int64 = 0
	var polylineOffset int64 = 0

	for _, gmNdID := range gmNodeIdsSorted {
		waysFrom := ways[gmNdID]

		var polylines []string
		var edgeToRecords []edgeToRecord

		for edgeTo, v := range waysFrom {
			polyline := fs.pointEncoder.Encode(v.Points)
			polylineLen := int32(len(polyline))

			edgeToRec := edgeToRecord{
				nodeTo:           int32(edgeTo),
				distance:         v.Distance,
				wayType:          v.WayType,
				elevFrom: v.ElevationInfo.From,
				elevTo:   v.ElevationInfo.To,
				polylinePosition: polylinePosition{length: polylineLen, offset: polylineOffset},
			}
			edgeToRecords = append(edgeToRecords, edgeToRec)

			polylineOffset += int64(polylineLen)

			polylines = append(polylines, polyline)
		}
		err := fs.writePolylinesFile(polylines)
		if err != nil {
			return err
		}

		err = fs.writeEdgeToFile(edgeToRecords)
		if err != nil {
			return err
		}

		edgeStart := edgeStartRecord{int32(len(waysFrom)), offset}
		err = fs.writeEdgeFromFile(gmNdID, edgeStart)
		if err != nil {
			return err
		}

		offset += int64(len(waysFrom)) * edgeToIndividualRecordSize
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