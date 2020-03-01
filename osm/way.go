package osm

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
)

type osmFileRead struct {
	osmFilename string
	nodeDB gravelmap.Osm2GmNodeReaderWriter
	gmNodeRd gravelmap.GmNodeReader
	wayStorer gravelmap.WayStorer
	graphWayAdder gravelmap.GraphWayAdder
	costEvaluator gravelmap.CostEvaluator
}

func NewOsmWayFileRead(
	osmFilename string,
	nodeDB gravelmap.Osm2GmNodeReaderWriter,
	gmNodeRd gravelmap.GmNodeReader,
	wayStorer gravelmap.WayStorer,
	graphWayAdder gravelmap.GraphWayAdder,
	costEvaluator gravelmap.CostEvaluator,
	) *osmFileRead {
	return &osmFileRead{
		osmFilename: osmFilename,
		nodeDB: nodeDB,
		gmNodeRd: gmNodeRd,
		wayStorer: wayStorer,
		graphWayAdder: graphWayAdder,
		costEvaluator: costEvaluator,
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

	// TODO: inject this to the constructor
	osm2GmWays := NewOsm2GmWays(fs.nodeDB, fs.gmNodeRd, fs.costEvaluator)

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				osm2GmWays.Add(v.NodeIDs, v.Tags)
			default:
				break
			}
		}
	}

	fs.graphWayAdder.AddWays(osm2GmWays.Get())
	return fs.wayStorer.Store(osm2GmWays.Get())
}
