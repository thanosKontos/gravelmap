package way

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
	"googlemaps.github.io/maps"
)

type fileStore struct {
	destinationDir string
	osmFilename string
	nodeDB gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd gravelmap.GmNodeReader
}

func NewWayFileStore(destinationDir string, osmFilename string, nodeDB gravelmap.Osm2GmNodeReaderWriter, gmNodeRd gravelmap.GmNodeReader) *fileStore {
	return &fileStore{
		destinationDir: destinationDir,
		osmFilename: osmFilename,
		nodeDB: nodeDB,
		gmNodeRd: gmNodeRd,
	}
}

type wayTo struct {
	ndTo int
	polyline string
}

func (fs *fileStore) Persist() error {
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

	ways := make(map[int][]wayTo)
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				prevEdge := 0
				var wayNd []int
				for i, nd := range v.NodeIDs {
					osm2gm := fs.nodeDB.Read(nd)
					wayNd = append(wayNd, osm2gm.NewID)

					if i == 0 {
						prevEdge = fs.nodeDB.Read(nd).NewID
					} else if i == len(v.NodeIDs) - 1 {
						gmID := fs.nodeDB.Read(nd).NewID

						ways[gmID] = append(ways[gmID], wayTo{prevEdge, fs.getWayPolyline(wayNd)})
						ways[prevEdge] = append(ways[prevEdge], wayTo{gmID, fs.getWayPolyline(wayNd)})

						wayNd = []int{prevEdge}
					} else {
						gmNd := fs.nodeDB.Read(nd)
						if gmNd.Occurrences > 1 {

							ways[gmNd.NewID] = append(ways[gmNd.NewID], wayTo{prevEdge, fs.getWayPolyline(wayNd)})
							ways[prevEdge] = append(ways[prevEdge], wayTo{gmNd.NewID, fs.getWayPolyline(wayNd)})

							prevEdge = gmNd.NewID
							wayNd = []int{prevEdge}
						}
					}
				}
			default:
				break
			}
		}
	}

	return fs.writeFilesForWays(ways)
}

func (fs *fileStore) getWayPolyline(wayNds []int) string {
	var latLngs []maps.LatLng
	for _, gmNdID := range wayNds {
		gmNode, _ := fs.gmNodeRd.Read(int32(gmNdID))
		latLngs = append(latLngs, maps.LatLng{Lat: gmNode.Lat, Lng: gmNode.Lng})
	}

	return maps.Encode(latLngs)
}

type nodeStartRecord struct {
	connectionsCnt int32
	nodeToOffset int64
}

func (fs *fileStore) writeFilesForWays(ways map[int][]wayTo) error {
	f, err := os.Create(fmt.Sprintf("%s/%s", fs.destinationDir, "node_start.bin"))
	defer f.Close()
	if err != nil {
		return err
	}

	recordSize := 8

	var offset int64 = 0
	for gmNdID, way := range ways {
		f.Seek(int64(gmNdID*recordSize), 0)
		nodeStart := nodeStartRecord{int32(len(way)), offset}

		var buf bytes.Buffer
		err := binary.Write(&buf, binary.BigEndian, nodeStart)
		if err != nil {
			return err
		}

		n, err := f.Write(buf.Bytes())
		if err != nil {
			return err
		}

		offset += int64(n)
 	}

	return nil
}

