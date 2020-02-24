package prepare

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/dijkstra"
)

type graph struct {
	osmFilename string
	graph *dijkstra.Graph
	nodeDB gravelmap.Osm2GmNodeReaderWriter
	distanceCalc gravelmap.DistanceCalculator
}

func NewGraph(osmFilename string, nodeDB gravelmap.Osm2GmNodeReaderWriter, distanceCalc gravelmap.DistanceCalculator) *graph {
	return &graph{
		osmFilename: osmFilename,
		graph: dijkstra.NewGraph(),
		nodeDB: nodeDB,
		distanceCalc: distanceCalc,
	}
}

func (g *graph) Prepare () {
	f, err := os.Open(g.osmFilename)
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

	var lastAddedVertex = 0
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				var gmNds []gravelmap.NodeOsm2GM
				for _, osmNdID := range v.NodeIDs {
					gmNode := g.nodeDB.Read(osmNdID)
					gmNds = append(gmNds, *gmNode)
				}

				vtx := g.addWaysWithCostToGraph(gmNds, lastAddedVertex)
				if vtx != -1 {
					lastAddedVertex = vtx
				}
			}
		}
	}
}

// TODO this is incorrect. Need to create an abstraction in GM domain and return this instead
// But will leave this technical debt for the POC
func (g *graph) GetGraph () *dijkstra.Graph {
	return g.graph
}

func (g *graph) addWaysWithCostToGraph(wayNds []gravelmap.NodeOsm2GM, previousLastAddedVertex int) int {
	var previousEdge gravelmap.NodeOsm2GM
	var firstEdge gravelmap.NodeOsm2GM
	lastAddedVertex := -1

	for _, wayNd := range wayNds {
		if isEdge := wayNd.Occurrences > 1; isEdge {
			if wayNd.NewID > previousLastAddedVertex || previousLastAddedVertex == 0 {
				g.graph.AddVertex(wayNd.NewID)
				lastAddedVertex = wayNd.NewID
			}

			if isFirstEdge := firstEdge == (gravelmap.NodeOsm2GM{}); isFirstEdge {
				firstEdge = wayNd
				previousEdge = wayNd
				continue
			}



			g.graph.AddArc(wayNd.NewID, previousEdge.NewID, 1)
			g.graph.AddArc(previousEdge.NewID, wayNd.NewID, 1)

			previousEdge = wayNd
		}

	}

	return lastAddedVertex
}
