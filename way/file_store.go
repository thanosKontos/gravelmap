package way

import (
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
	nodeDB gravelmap.NodeOsm2GMReaderWriter
	gmNodeRd gravelmap.GmNodeReader
}

func NewWayFileStore(destinationDir string, osmFilename string, nodeDB gravelmap.NodeOsm2GMReaderWriter, gmNodeRd gravelmap.GmNodeReader) *fileStore {
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

				//fmt.Println("\n============")
			default:
				break
			}
		}
	}

	return nil
}

func (fs *fileStore) getWayPolyline(wayNds []int) string {
	var latLngs []maps.LatLng
	for _, nd := range wayNds {
		gmNode, _ := fs.gmNodeRd.Read(int32(nd))
		latLngs = append(latLngs, maps.LatLng{Lat: gmNode.Lat, Lng: gmNode.Lng})
	}

	return maps.Encode(latLngs)
}
