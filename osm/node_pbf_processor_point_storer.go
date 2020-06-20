package osm

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
)

type osmNodeProcess struct {
	osmFilename      string
	osm2GmStore      gravelmap.Osm2GmNodeReaderWriter
	NodePointStorer  gravelmap.NodePointStorer
	osm2LatLngWriter gravelmap.Osm2LatLngWriter
}

func NewOsmNodeProcessor(
	osmFilename string,
	osm2GmStore gravelmap.Osm2GmNodeReaderWriter,
	NodePointStorer gravelmap.NodePointStorer,
	osm2LatLngWriter gravelmap.Osm2LatLngWriter,
) *osmNodeProcess {
	return &osmNodeProcess{
		osmFilename:      osmFilename,
		osm2GmStore:      osm2GmStore,
		NodePointStorer:  NodePointStorer,
		osm2LatLngWriter: osm2LatLngWriter,
	}
}

func (fs *osmNodeProcess) Process() error {
	f, err := os.Open(fs.osmFilename)
	if err != nil {
		return err
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		return err
	}

	var nodePtBatch []gravelmap.NodePoint

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

				nodePt := gravelmap.NodePoint{NodeID: int32(gm2OsmNode.ID), Point: gravelmap.Point{Lat: v.Lat, Lng: v.Lon}}

				// Write nodes in file in order to be able to find lat long per id
				fs.osm2LatLngWriter.Write(int(nodePt.NodeID), nodePt.Point)

				// Write edge in bounding boxes in order to be able to find closest edge per lat/lng
				if gm2OsmNode.ConnectionCnt > 1 {
					nodePtBatch = append(nodePtBatch, nodePt)

					if len(nodePtBatch) >= 10000 {
						fs.NodePointStorer.BatchStore(nodePtBatch)
						nodePtBatch = []gravelmap.NodePoint{}
					}
				}

			default:
				break
			}
		}
	}

	if len(nodePtBatch) > 0 {
		fs.NodePointStorer.BatchStore(nodePtBatch)
	}

	return nil
}
