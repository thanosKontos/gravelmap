package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
)

const (
	recordSize = 16
	filename = "node.bin"
)

type fileStore struct {
	destinationDir string
	osmFilename string
	nodeDB gravelmap.Osm2GmNodeReaderWriter
}

func NewNodeFileStore(destinationDir string, osmFilename string, nodeDB gravelmap.Osm2GmNodeReaderWriter) *fileStore {
	return &fileStore{
		destinationDir: destinationDir,
		osmFilename: osmFilename,
		nodeDB: nodeDB,
	}
}

func (fs *fileStore) Persist() error {
	f, err := os.Create(fmt.Sprintf("%s/%s", fs.destinationDir, filename))
	defer f.Close()
	if err != nil {
		return err
	}

	osmF, err := os.Open(fs.osmFilename)
	if err != nil {
		return err
	}
	defer osmF.Close()

	d := osmpbf.NewDecoder(osmF)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return err
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				gm2OsmNode := fs.nodeDB.Read(v.ID)
				if gm2OsmNode == nil {
					continue
				}

				gmNd := gravelmap.GMNode{ID: int32(gm2OsmNode.NewID), Point: gravelmap.Point{Lat: v.Lat, Lng: v.Lon}}

				writeGmNode(f, gmNd)
			default:
				break
			}
		}
	}

	return nil
}

func writeGmNode(f *os.File, gmNd gravelmap.GMNode) {
	_, err := f.Seek(int64(gmNd.ID*recordSize), 0)
	if err != nil {
		fmt.Println(err)
	}

	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, gmNd.Point)
	if err != nil {
		fmt.Println(err)
	}
	writeNextBytes(f, buf.Bytes())
}
func writeNextBytes(file *os.File, bytes []byte) {
	_, err := file.Write(bytes)

	if err != nil {
		log.Fatal(err)
	}
}

func (fs *fileStore) Read(ndID int32) (*gravelmap.GMNode, error) {
	f, err := os.Open(fmt.Sprintf("%s/%s", fs.destinationDir, filename))
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var pt gravelmap.Point

	f.Seek(int64(ndID*recordSize), 0)
	data := readNextBytes(f, recordSize)
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.BigEndian, &pt)
	if err != nil {
		log.Fatal("binary.Read failed", err)
	}

	return &gravelmap.GMNode{ID: ndID, Point: pt}, nil
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}