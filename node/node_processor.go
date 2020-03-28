package node

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
	edgeBatchStorer  gravelmap.EdgeBatchStorer
	osm2LatLngWriter gravelmap.Osm2LatLngWriter
}

func NewOsmNodeProcessor(
	osmFilename string,
	osm2GmStore gravelmap.Osm2GmNodeReaderWriter,
	edgeBatchStorer gravelmap.EdgeBatchStorer,
	osm2LatLngWriter gravelmap.Osm2LatLngWriter,
) *osmNodeProcess {
	return &osmNodeProcess{
		osmFilename:      osmFilename,
		osm2GmStore:      osm2GmStore,
		edgeBatchStorer:  edgeBatchStorer,
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

				gmNd := gravelmap.Node{Id: gm2OsmNode.Id, Point: gravelmap.Point{Lat: v.Lat, Lng: v.Lon}}

				// Write nodes in file in order to be able to find lat long per id
				fs.osm2LatLngWriter.Write(gmNd.Id, gmNd.Point)

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
