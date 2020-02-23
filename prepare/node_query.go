package prepare

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

func NewOsm2GmEdge(osmFilename string, osm2GmNodeRw gravelmap.Osm2GmNodeReaderWriter) *osm2GmEdge {
	return &osm2GmEdge{
		osmFilename:  osmFilename,
		osm2GmNodeRw: osm2GmNodeRw,
	}
}

func (n *osm2GmEdge) Extract() gravelmap.Osm2GmNodeReaderWriter {
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
		log.Fatal(err)
	}

	inc := 0
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				for _, osmNdID := range v.NodeIDs {
					ndDB := n.osm2GmNodeRw.Read(osmNdID)

					if ndDB == nil {
						inc++
						_ = n.osm2GmNodeRw.Write(&gravelmap.NodeOsm2GM{OldID: osmNdID, NewID: inc, Occurrences: 1})
					} else {
						newCnt := ndDB.Occurrences + 1
						_ = n.osm2GmNodeRw.Write(&gravelmap.NodeOsm2GM{OldID: ndDB.OldID, NewID: ndDB.NewID, Occurrences: newCnt})
					}
				}
			}
		}
	}

	return n.osm2GmNodeRw
}
