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
	costEvaluator gravelmap.CostEvaluator
}

func NewGraph(osmFilename string, nodeDB gravelmap.Osm2GmNodeReaderWriter, costEvaluator gravelmap.CostEvaluator) *graph {
	return &graph{
		osmFilename: osmFilename,
		graph: dijkstra.NewGraph(),
		nodeDB: nodeDB,
		costEvaluator: costEvaluator,
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

				vtx := g.addWaysWithCostToGraph(gmNds, v.Tags, lastAddedVertex)
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

func (g *graph) addWaysWithCostToGraph(wayNdsOsm2GM []gravelmap.NodeOsm2GM, tags map[string]string, previousLastAddedVertex int) int {
	var evaluativeWay = gravelmap.EvaluativeWay{Tags: tags}
	var previousSubwayPoint = gravelmap.Point{}
	var previousEdge gravelmap.NodeOsm2GM
	var firstEdge gravelmap.NodeOsm2GM
	lastAddedVertex := -1

	for _, ndOsm2GM := range wayNdsOsm2GM {
		if isEdge := ndOsm2GM.Occurrences > 1; isEdge {
			if ndOsm2GM.GmID > previousLastAddedVertex || previousLastAddedVertex == 0 {
				g.graph.AddVertex(ndOsm2GM.GmID)
				lastAddedVertex = ndOsm2GM.GmID
			}

			if isFirstEdge := firstEdge == (gravelmap.NodeOsm2GM{}); isFirstEdge {
				evaluativeWay.Points = append(evaluativeWay.Points, ndOsm2GM.Point)
				previousSubwayPoint = ndOsm2GM.Point
				firstEdge = ndOsm2GM
				previousEdge = ndOsm2GM
				continue
			}

			evaluativeWay.Points = append(evaluativeWay.Points, ndOsm2GM.Point)

			cost := g.costEvaluator.Evaluate(evaluativeWay)

			g.graph.AddArc(ndOsm2GM.GmID, previousEdge.GmID, cost.Cost)
			g.graph.AddArc(previousEdge.GmID, ndOsm2GM.GmID, cost.ReverseCost)

			evaluativeWay.Points = []gravelmap.Point{ndOsm2GM.Point}

			previousEdge = ndOsm2GM
			previousSubwayPoint = ndOsm2GM.Point
		} else {
			if hasPreviousSubwayPoint := previousSubwayPoint != (gravelmap.Point{}); hasPreviousSubwayPoint {
				evaluativeWay.Points = append(evaluativeWay.Points, ndOsm2GM.Point)
			}
		}
	}

	return lastAddedVertex
}
