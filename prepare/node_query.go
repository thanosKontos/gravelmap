package prepare

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
)

type nodeQuery struct {
	osmFilename  string
	ndDB gravelmap.NodeOsm2GMReaderWriter
}

func NewNodeQuerer(osmFilename string, ndDB gravelmap.NodeOsm2GMReaderWriter) *nodeQuery {
	return &nodeQuery{
		osmFilename:  osmFilename,
		ndDB: ndDB,
	}
}

func (n *nodeQuery) Prepare () gravelmap.NodeOsm2GMReaderWriter {
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
					ndDB := n.ndDB.Read(osmNdID)

					if ndDB == nil {
						inc++
						n.ndDB.Write(&gravelmap.NodeOsm2GM{osmNdID, inc, 1})
					} else {
						newCnt := ndDB.Occurrences + 1
						n.ndDB.Write(&gravelmap.NodeOsm2GM{ndDB.OldID, ndDB.NewID, newCnt})
					}
				}
			}
		}
	}

	return n.ndDB
}
