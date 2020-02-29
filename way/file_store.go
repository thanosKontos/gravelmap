package way

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
	"googlemaps.github.io/maps"
)

type fileStore struct {
	storageDir string
	osmFilename string
	nodeDB gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd gravelmap.GmNodeReader
}

func NewWayFileStore(storageDir string, osmFilename string, nodeDB gravelmap.Osm2GmNodeReaderWriter, gmNodeRd gravelmap.GmNodeReader) *fileStore {
	return &fileStore{
		storageDir: storageDir,
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

	wayNds := make(map[int][]wayTo)
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
						prevEdge = fs.nodeDB.Read(osmNdID).GmID
					} else if i == len(v.NodeIDs) - 1 {
						gmID := fs.nodeDB.Read(osmNdID).GmID

						wayNds[gmID] = append(wayNds[gmID], wayTo{prevEdge, fs.getWayPolyline(wayGmNds, true)})
						wayNds[prevEdge] = append(wayNds[prevEdge], wayTo{gmID, fs.getWayPolyline(wayGmNds, false)})

						wayGmNds = []int{prevEdge}
					} else {
						if gmNd := fs.nodeDB.Read(osmNdID); gmNd.Occurrences > 1 {
							wayNds[gmNd.GmID] = append(wayNds[gmNd.GmID], wayTo{prevEdge, fs.getWayPolyline(wayGmNds, true)})
							wayNds[prevEdge] = append(wayNds[prevEdge], wayTo{gmNd.GmID, fs.getWayPolyline(wayGmNds, false)})

							prevEdge = gmNd.GmID
							wayGmNds = []int{prevEdge}
						}
					}
				}
			default:
				break
			}
		}
	}

	return fs.writeFilesForWays(wayNds)
}

func (fs *fileStore) getWayPolyline(wayGmNds []int, reverse bool) string {
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

type polylineLookupRecord struct {
	nodeTo int32
	polylineLen int32
	polylineOffset int64
}

func (fs *fileStore) writeFilesForWays(ways map[int][]wayTo) error {
	var gmNodeIdsSorted []int
	for k := range ways {
		gmNodeIdsSorted = append(gmNodeIdsSorted, k)
	}
	sort.Ints(gmNodeIdsSorted)

	esFl, err := os.Create(fmt.Sprintf("%s/%s", fs.storageDir, edgeStartFilename))
	defer esFl.Close()
	if err != nil {
		return err
	}

	plLkFl, err := os.Create(fmt.Sprintf("%s/%s", fs.storageDir, edgeToFilename))
	defer plLkFl.Close()
	if err != nil {
		return err
	}

	plFl, err := os.Create(fmt.Sprintf("%s/%s", fs.storageDir, polylinesFilename))
	defer plFl.Close()
	if err != nil {
		return err
	}

	var offset int64 = 0
	var polylineOffset int64 = 0

	for _, gmNdID := range gmNodeIdsSorted {
		way := ways[gmNdID]

		var polylines []string
		var wayToPolylineLookups []polylineLookupRecord

		for i := 0; i < len(way); i++ {
			polylineLen := int32(len(way[i].polyline))

			polylineLookupRec := polylineLookupRecord{int32(way[i].ndTo), polylineLen, polylineOffset}
			wayToPolylineLookups = append(wayToPolylineLookups, polylineLookupRec)

			polylineOffset += int64(polylineLen)

			polylines = append(polylines, way[i].polyline)
		}
		err = fs.writePolylinesFile(plFl, polylines)
		if err != nil {
			return err
		}

		err = fs.writePolylinesLookupFile(plLkFl, wayToPolylineLookups)
		if err != nil {
			return err
		}

		edgeStart := edgeStartRecord{int32(len(way)), offset}

		err = fs.writeEdgeFromFile(esFl, gmNdID, edgeStart)
		if err != nil {
			return err
		}

		offset += 4*int64(len(way))*2+int64(len(way))*8
	}

	return nil
}

func (fs *fileStore) writePolylinesLookupFile(f *os.File, plsLookup []polylineLookupRecord) error {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, plsLookup)
	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())
	return err
}

func (fs *fileStore) writePolylinesFile(f *os.File, pls []string) error {
	for _, pl := range pls {
		_, err := f.WriteString(pl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *fileStore) writeEdgeFromFile(f *os.File, edgeStartId int, edgeStart edgeStartRecord) error {
	f.Seek(int64(edgeStartId*edgeStartRecordSize), 0)

	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, edgeStart)
	if err != nil {
		return err
	}

	_, err = f.Write(buf.Bytes())

	return err
}