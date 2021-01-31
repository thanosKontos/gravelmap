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
	osm2gmRecordSize = 22
	osm2gmFilename   = "osm2gm.bin"
)

type connNodeFile struct {
	ID      int32           // 4
	ConnCnt int16           // 2
	Point   gravelmap.Point // 16
}

type osm2GmNodeFile struct {
	file *os.File
}

func NewOsm2GmNodeFileStore(destinationDir string) *osm2GmNodeFile {
	f, _ := os.Create(fmt.Sprintf("%s/%s", destinationDir, osm2gmFilename))
	return &osm2GmNodeFile{
		file: f,
	}
}

func (nf osm2GmNodeFile) Write(osmNdID int64, gm *gravelmap.ConnectionNode) error {
	if _, err := nf.file.Seek(int64(osmNdID*osm2gmRecordSize), 0); err != nil {
		fmt.Println(err)
	}

	connNodeFile := connNodeFile{int32(gm.ID), int16(gm.ConnectionCnt), gm.Point}

	var buf bytes.Buffer
	if err := binary.Write(&buf, binary.BigEndian, connNodeFile); err != nil {
		fmt.Println(err)
	}
	writeNextBytes(nf.file, buf.Bytes())

	return nil
}

func (nf osm2GmNodeFile) Read(osmNdID int64) *gravelmap.ConnectionNode {
	var cnf connNodeFile

	nf.file.Seek(int64(osmNdID*osm2gmRecordSize), 0)
	data := readNextBytes(nf.file, osm2gmRecordSize)
	buffer := bytes.NewBuffer(data)
	if err := binary.Read(buffer, binary.BigEndian, &cnf); err != nil {
		log.Fatal("binary.Read failed", err)
	}

	if cnf.ID == 0 {
		return nil
	}

	cn := gravelmap.ConnectionNode{int(cnf.ID), int(cnf.ConnCnt), cnf.Point}

	return &cn
}
