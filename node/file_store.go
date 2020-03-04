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
	filename   = "node_positions.bin"
)

type fileStore struct {
	destinationDir  string
	osmFilename     string
	osm2GmStore     gravelmap.Osm2GmNodeReaderWriter
	edgeBatchStorer gravelmap.EdgeBatchStorer
}

func NewNodeFileStore(
	destinationDir string,
	osmFilename string,
	osm2GmStore gravelmap.Osm2GmNodeReaderWriter,
	edgeBatchStorer gravelmap.EdgeBatchStorer,
) *fileStore {
	return &fileStore{
		destinationDir:  destinationDir,
		osmFilename:     osmFilename,
		osm2GmStore:     osm2GmStore,
		edgeBatchStorer: edgeBatchStorer,
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

	var gmNdBatch []gravelmap.Node

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				gm2OsmNode := fs.osm2GmStore.Read(v.ID)
				if gm2OsmNode == nil {
					continue
				}

				gm2OsmNode.Point = gravelmap.Point{Lat: v.Lat, Lng: v.Lon}
				_ = fs.osm2GmStore.Write(v.ID, gm2OsmNode)

				// TODO: create an extract node service and create a node package to include the 2 jobs below
				// inject to the service a osmFilename, nodePositionWriter, gmEdgeBboxWriter (the implementation will be file)

				gmNd := gravelmap.Node{Id: gm2OsmNode.Id, Point: gravelmap.Point{Lat: v.Lat, Lng: v.Lon}}

				// Write nodes in file in order to be able to find lat long per id
				writeGmNode(f, gmNd)

				// Write edge in bounding boxes in order to be able to find closest edge per lat/lng
				if gm2OsmNode.Occurrences > 1 {
					gmNdBatch = append(gmNdBatch, gmNd)

					if len(gmNdBatch) >= 10000 {
						fs.edgeBatchStorer.BatchStore(gmNdBatch)
						gmNdBatch = []gravelmap.Node{}
					}
				}

			default:
				break
			}
		}
	}

	if len(gmNdBatch) > 0 {
		fs.edgeBatchStorer.BatchStore(gmNdBatch)
	}

	return nil
}

func writeGmNode(f *os.File, gmNd gravelmap.Node) {
	_, err := f.Seek(int64(gmNd.Id*recordSize), 0)
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

func (fs *fileStore) Read(ndID int) (*gravelmap.Node, error) {
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

	return &gravelmap.Node{Id: ndID, Point: pt}, nil
}

func readNextBytes(f *os.File, number int) []byte {
	byteSeq := make([]byte, number)
	_, _ = f.Read(byteSeq)

	return byteSeq
}
