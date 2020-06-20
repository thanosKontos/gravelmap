package osm

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
)

type osm2GmEdge struct {
	osmFilename  string
	osm2GmNodeRw gravelmap.Osm2GmNodeReaderWriter
}

func NewOsmWayProcessor(osmFilename string, osm2GmNodeRw gravelmap.Osm2GmNodeReaderWriter) *osm2GmEdge {
	return &osm2GmEdge{
		osmFilename:  osmFilename,
		osm2GmNodeRw: osm2GmNodeRw,
	}
}

func (n *osm2GmEdge) Process() error {
	f, err := os.Open(n.osmFilename)
	if err != nil {
		log.Fatal(err)
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

	inc := 0
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return err
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				for _, osmNdID := range v.NodeIDs {
					ndDB := n.osm2GmNodeRw.Read(osmNdID)

					if ndDB == nil {
						inc++
						_ = n.osm2GmNodeRw.Write(osmNdID, &gravelmap.ConnectionNode{ID: inc, ConnectionCnt: 1})
					} else {
						newCnt := ndDB.ConnectionCnt + 1
						_ = n.osm2GmNodeRw.Write(osmNdID, &gravelmap.ConnectionNode{ID: ndDB.ID, ConnectionCnt: newCnt})
					}
				}
			}
		}
	}

	return nil
}
