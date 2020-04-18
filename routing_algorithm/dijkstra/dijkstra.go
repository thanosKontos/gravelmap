package dijkstra

import (
	"math"

	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/graph"
)

type dijkstraRouter struct {
	graph *graph.Graph

	costToDest int64
	destFound  bool
	list       graph.DijkstraList
}

func NewDijkstraRouter(graph *graph.Graph) *dijkstraRouter {
	return &dijkstraRouter{
		graph: graph,
	}
}

//Shortest calculates the shortest path from src to dest
func (g *dijkstraRouter) FindShortest(src, dest int) (gravelmap.BestPath, error) {
	g.setup(src)
	return g.postSetupEvaluate(src, dest)
}

func (g *dijkstraRouter) setup(src int) {
	g.setupList()
	g.setDefaults()

	//Set the cost of initial vertex 0 and add it to the list
	g.graph.Vertices[src].Cost = 0
	g.list.PushOrdered(&g.graph.Vertices[src])
}

func (g *dijkstraRouter) setupList() {
	if len(g.graph.Vertices) < 800 {
		g.list = graph.LinkedListNewLong()
		return
	}

	g.list = graph.PriorityQueueNewLong()
	return
}

func (g *dijkstraRouter) postSetupEvaluate(src, dest int) (gravelmap.BestPath, error) {
	var current *graph.Vertex
	oldCurrent := -1
	for g.list.Len() > 0 {
		//Visit the current lowest distanced Vertex
		//TODO WTF
		current = g.list.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		//If the current cost is already worse than the best one try another Vertex
		if current.Cost >= g.costToDest {
			continue
		}
		for v, dist := range current.Arcs {
			//If the arc has better access, than the current costToDest, update the Vertex being touched
			if current.Cost+dist < g.graph.Vertices[v].Cost {
				if current.BestVertices[0] == v && g.graph.Vertices[v].ID != dest {
					//also only do this if we aren't checkout out the best cost again
					//This seems familiar 8^)
					return gravelmap.BestPath{}, newErrLoop(current.ID, v)
				}
				g.graph.Vertices[v].Cost = current.Cost + dist
				g.graph.Vertices[v].BestVertices[0] = current.ID
				if v == dest {
					//If this is the destination update costToDest, so we can stop looking at useless Vertices
					g.costToDest = current.Cost + dist
					g.destFound = true
					continue // Do not push if dest
				}

				//Push this updated Vertex into the list to be evaluated, pushes in sorted form
				g.list.PushOrdered(&g.graph.Vertices[v])
			}
		}
	}
	return g.finally(src, dest)
}

func (g *dijkstraRouter) finally(src, dest int) (gravelmap.BestPath, error) {
	if !g.destFound {
		return gravelmap.BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}

func (g *dijkstraRouter) bestPath(src, dest int) gravelmap.BestPath {
	var path []int
	for c := g.graph.Vertices[dest]; c.ID != src; c = g.graph.Vertices[c.BestVertices[0]] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return gravelmap.BestPath{Distance: g.graph.Vertices[dest].Cost, Path: path}
}

// 1. Reset state
// 2. Reset the cost to destination
// set all best vertices to -1 (unused) and set the defaults *almost* as bad
func (g *dijkstraRouter) setDefaults() {
	g.destFound = false
	g.costToDest = int64(math.MaxInt64)

	for i := range g.graph.Vertices {
		g.graph.Vertices[i].BestVertices = []int{-1}
		g.graph.Vertices[i].Cost = int64(math.MaxInt64) - 2
	}
}
