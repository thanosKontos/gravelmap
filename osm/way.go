package osm

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
	"googlemaps.github.io/maps"
)

type osmFileRead struct {
	osmFilename string
	nodeDB gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd gravelmap.GmNodeReader
	wayStorer gravelmap.WayStorer
}

func NewOsmWayFileRead(
	osmFilename string,
	nodeDB gravelmap.Osm2GmNodeReaderWriter,
	gmNodeRd gravelmap.GmNodeReader,
	wayStorer gravelmap.WayStorer,
	) *osmFileRead {
	return &osmFileRead{
		osmFilename: osmFilename,
		nodeDB: nodeDB,
		gmNodeRd: gmNodeRd,
		wayStorer: wayStorer,
	}
}

func (fs *osmFileRead) Process() error {
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

	wayNds := make(map[int][]gravelmap.WayTo)
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				prevEdge := 0
				var wayGmNds []int
				for i, osmNdID := range v.NodeIDs {
					osm2gm := fs.nodeDB.Read(osmNdID)
					wayGmNds = append(wayGmNds, osm2gm.GmID)

					if i == 0 {
						prevEdge = osm2gm.GmID
					} else if i == len(v.NodeIDs) - 1 {
						wayNds[osm2gm.GmID] = append(wayNds[osm2gm.GmID], gravelmap.WayTo{prevEdge, fs.getWayPolyline(wayGmNds, true)})
						wayNds[prevEdge] = append(wayNds[prevEdge], gravelmap.WayTo{osm2gm.GmID, fs.getWayPolyline(wayGmNds, false)})

						wayGmNds = []int{prevEdge}
					} else {
						if osm2gm.Occurrences > 1 {
							wayNds[osm2gm.GmID] = append(wayNds[osm2gm.GmID], gravelmap.WayTo{prevEdge, fs.getWayPolyline(wayGmNds, true)})
							wayNds[prevEdge] = append(wayNds[prevEdge], gravelmap.WayTo{osm2gm.GmID, fs.getWayPolyline(wayGmNds, false)})

							prevEdge = osm2gm.GmID
							wayGmNds = []int{prevEdge}
						}
					}
				}
			default:
				break
			}
		}
	}

	return fs.wayStorer.Store(wayNds)
}

func (fs *osmFileRead) getWayPolyline(wayGmNds []int, reverse bool) string {
	var latLngs []maps.LatLng

	if reverse {
		for i := len(wayGmNds)-1; i >= 0; i-- {
			gmNode, _ := fs.gmNodeRd.Read(int32(wayGmNds[i]))
			latLngs = append(latLngs, maps.LatLng{Lat: gmNode.Lat, Lng: gmNode.Lng})
		}
	} else {
		for _, gmNdID := range wayGmNds {
			gmNode, _ := fs.gmNodeRd.Read(int32(gmNdID))
			latLngs = append(latLngs, maps.LatLng{Lat: gmNode.Lat, Lng: gmNode.Lng})
		}
	}

	return maps.Encode(latLngs)
}

