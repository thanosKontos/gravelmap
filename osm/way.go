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
	osmFilename    string
	wayStorer      gravelmap.WayStorer
	graphWayAdder  gravelmap.GraphWayAdder
	wayAdderGetter gravelmap.WayAdderGetter
}

func NewOsmWayFileRead(
	osmFilename string,
	wayStorer gravelmap.WayStorer,
	graphWayAdder gravelmap.GraphWayAdder,
	wayAdderGetter gravelmap.WayAdderGetter,
) *osmFileRead {
	return &osmFileRead{
		osmFilename:    osmFilename,
		wayStorer:      wayStorer,
		graphWayAdder:  graphWayAdder,
		wayAdderGetter: wayAdderGetter,
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

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				fs.wayAdderGetter.Add(v.NodeIDs, v.Tags)
			default:
				break
			}
		}
	}

	fs.graphWayAdder.AddWays(fs.wayAdderGetter.Get())
	return fs.wayStorer.Store(fs.wayAdderGetter.Get())
}
