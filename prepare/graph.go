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
	nodeDB gravelmap.NodeOsm2GMReaderWriter
}

func NewGraph(osmFilename string, nodeDB gravelmap.NodeOsm2GMReaderWriter) *graph {
	return &graph{
		osmFilename: osmFilename,
		graph: dijkstra.NewGraph(),
		nodeDB: nodeDB,
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
				var intersections []int
				for _, osmNdID := range v.NodeIDs {
					edge := g.getEdge(osmNdID)
					if edge != nil {
						intersections = append(intersections, edge.NewID)
					}
				}

				vtx := g.addIntersectionsToGraph(intersections, lastAddedVertex)
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

func (g *graph) getEdge (nd int64) *gravelmap.NodeOsm2GM {
	node := g.nodeDB.Read(nd)
	if node.Occurrences > 1 {
		return node
	}

	return nil
}

func (g *graph) addIntersectionsToGraph(intersections []int, previousLastAddedVertex int) int {
	previous := 0
	lastAddedVertex := -1

	for i, isn := range intersections {
		if isn > previousLastAddedVertex || previousLastAddedVertex == 0 {
			g.graph.AddVertex(isn)
			lastAddedVertex = isn
		}

		if i == 0 {
			previous = isn
		} else {
			g.graph.AddArc(isn, previous, 1)
			g.graph.AddArc(previous, isn, 1)
		}
	}

	return lastAddedVertex
}
