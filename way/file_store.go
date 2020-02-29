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
}

func NewFileStore(storageDir string) *fileStore {
	return &fileStore{
		storageDir: storageDir,
	}
}

type polylineLookupRecord struct {
	nodeTo int32
	polylineLen int32
	polylineOffset int64
}

func (fs *fileStore) Store(ways map[int][]gravelmap.WayTo) error {
	var gmNodeIdsSorted []int
	for k := range ways {
		gmNodeIdsSorted = append(gmNodeIdsSorted, k)
	}
	sort.Ints(gmNodeIdsSorted)

	esFl, err := os.Create(fmt.Sprintf("%s/%s", fs.storageDir, edgeStartFilename))
	defer esFl.Close()
	if err != nil {
		return err
	}

	plLkFl, err := os.Create(fmt.Sprintf("%s/%s", fs.storageDir, edgeToFilename))
	defer plLkFl.Close()
	if err != nil {
		return err
	}

	plFl, err := os.Create(fmt.Sprintf("%s/%s", fs.storageDir, polylinesFilename))
	defer plFl.Close()
	if err != nil {
		return err
	}

	var offset int64 = 0
	var polylineOffset int64 = 0

	for _, gmNdID := range gmNodeIdsSorted {
		way := ways[gmNdID]

		var polylines []string
		var wayToPolylineLookups []polylineLookupRecord

		for i := 0; i < len(way); i++ {
			polylineLen := int32(len(way[i].Polyline))

			polylineLookupRec := polylineLookupRecord{int32(way[i].NdTo), polylineLen, polylineOffset}
			wayToPolylineLookups = append(wayToPolylineLookups, polylineLookupRec)

			polylineOffset += int64(polylineLen)

			polylines = append(polylines, way[i].Polyline)
		}
		err = fs.writePolylinesFile(plFl, polylines)
		if err != nil {
			return err
		}

		err = fs.writePolylinesLookupFile(plLkFl, wayToPolylineLookups)
		if err != nil {
			return err
		}

		edgeStart := edgeStartRecord{int32(len(way)), offset}

		err = fs.writeEdgeFromFile(esFl, gmNdID, edgeStart)
		if err != nil {
			return err
		}

		offset += 4*int64(len(way))*2+int64(len(way))*8
	}

	return nil
}

func (fs *fileStore) writePolylinesLookupFile(f *os.File, plsLookup []polylineLookupRecord) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, plsLookup)
	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())
	return err
}

func (fs *fileStore) writePolylinesFile(f *os.File, pls []string) error {
	for _, pl := range pls {
		_, err := f.WriteString(pl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *fileStore) writeEdgeFromFile(f *os.File, edgeStartId int, edgeStart edgeStartRecord) error {
	f.Seek(int64(edgeStartId*edgeStartRecordSize), 0)

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, edgeStart)
	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())

	return err
}
