package graph

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/thanosKontos/gravelmap"
)

const (
	filename = "graph_edges.bin"
)

type record struct {
	edgeFrom int32
	edgeTo   int32
	cost     int64
}

type fileStore struct {
	file *os.File
}

func NewGraphEdgeFileStore(destinationDir string) *fileStore {
	f, _ := os.Create(fmt.Sprintf("%s/%s", destinationDir, filename))
	return &fileStore{
		file: f,
	}
}

func (fs *fileStore) Store(ways map[int]map[int]gravelmap.Way) error {
	for edgeFromId, edgeFromWays := range ways {
		for edgeToId, way := range edgeFromWays {
			rec := record{int32(edgeFromId), int32(edgeToId), way.Cost}
			err := fs.writeGraphEdge(rec)
			if err != nil {
				return err
			}
		}
	}

	fs.file.Close()

	return nil
}

func (fs *fileStore) writeGraphEdge(rec record) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, rec)
	if err != nil {
		return err
	}

	_, err = fs.file.Write(buf.Bytes())
	return err
}
