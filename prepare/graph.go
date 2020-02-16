package prepare

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/go-memdb"
	"github.com/qedus/osmpbf"
	"github.com/thanosKontos/gravelmap/dijkstra"
)

type graph struct {
	osmFilename string
	graph *dijkstra.Graph
	nodeDB *memdb.MemDB
}

func NewGraph(osmFilename string, nodeDB *memdb.MemDB) *graph {
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
//fmt.Println(intersections)
				vtx := g.addIntersectionsToGraph(intersections, lastAddedVertex)
				if vtx != -1 {
					lastAddedVertex = vtx
				}
			}
		}
	}

	//best, err := g.graph.Shortest(2173, 2201)
	//best, err := g.graph.Shortest(1, 2)
	//best, err := g.graph.Shortest(214768, 214762)
	best, err := g.graph.Shortest(214425, 2)

	fmt.Println("Shortest distance ", best.Distance, " following path ", best.Path)
}

func (g *graph) getEdge (nd int64) *Node {
	rdTxn := g.nodeDB.Txn(false)

	raw, _ := rdTxn.First(nodeTable, "id", nd)
	rdTxn.Abort()

	if raw != nil {
		if raw.(*Node).Occurences > 1 {
			return raw.(*Node)
		}
	}

	return nil
}

func (g *graph) addIntersectionsToGraph(intersections []int, previousLastAddedVertex int) int {
	previous := 0
	lastAddedVertex := -1

	for i, isn := range intersections {
		if isn > previousLastAddedVertex || previousLastAddedVertex == 0 {
			g.graph.AddVertex(isn)
			//fmt.Println(isn)
			lastAddedVertex = isn
		}

		if i == 0 {
			previous = isn
		} else {
			//fmt.Println(isn, previous)
			g.graph.AddArc(isn, previous, 1)
			g.graph.AddArc(previous, isn, 1)
		}
	}

	//fmt.Println(lastAddedVertex)

	return lastAddedVertex
}
