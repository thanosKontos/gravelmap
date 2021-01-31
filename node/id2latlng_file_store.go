package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/thanosKontos/gravelmap"
)

const (
	recordSize = 16
	filename   = "node_positions.bin"
)

type fileStore struct {
	file *os.File
}

func NewOsm2LatLngFileStore(destinationDir string) *fileStore {
	f, _ := os.Create(fmt.Sprintf("%s/%s", destinationDir, filename))
	return &fileStore{
		file: f,
	}
}

func (fs *fileStore) Write(osmID int, point gravelmap.Point) {
	if _, err := fs.file.Seek(int64(osmID*recordSize), 0); err != nil {
		fmt.Println(err)
	}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, point); err != nil {
		fmt.Println(err)
	}
	writeNextBytes(fs.file, buf.Bytes())
}

func writeNextBytes(file *os.File, bytes []byte) {
	if _, err := file.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

func (fs *fileStore) Read(ndID int) (gravelmap.Point, error) {
	var pt gravelmap.Point

	fs.file.Seek(int64(ndID*recordSize), 0)
	data := readNextBytes(fs.file, recordSize)
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.BigEndian, &pt); err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return pt, nil
}

func readNextBytes(f *os.File, number int) []byte {
	byteSeq := make([]byte, number)
	_, _ = f.Read(byteSeq)

	return byteSeq
}

func (fs *fileStore) Close() {
	fs.file.Close()
}
