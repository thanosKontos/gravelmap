package dijkstra

import (
	"math"

	"github.com/thanosKontos/gravelmap"
)

type dijkstraRouter struct {
	graph      *Graph

	costToDest int64
	destFound  bool
	list       dijkstraList
}

func NewDijkstraRouter(graph *Graph) *dijkstraRouter {
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
	g.graph.Vertices[src].cost = 0
	g.list.PushOrdered(&g.graph.Vertices[src])
}

func (g *dijkstraRouter) setupList() {
	if len(g.graph.Vertices) < 800 {
		g.list = linkedListNewLong()
		return
	}

	g.list = priorityQueueNewLong()
	return
}

func (g *dijkstraRouter) postSetupEvaluate(src, dest int) (gravelmap.BestPath, error) {
	var current *Vertex
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
		if current.cost >= g.costToDest {
			continue
		}
		for v, dist := range current.Arcs {
			//If the arc has better access, than the current costToDest, update the Vertex being touched
			if current.cost+dist < g.graph.Vertices[v].cost {
				if current.bestVertices[0] == v && g.graph.Vertices[v].ID != dest {
					//also only do this if we aren't checkout out the best cost again
					//This seems familiar 8^)
					return gravelmap.BestPath{}, newErrLoop(current.ID, v)
				}
				g.graph.Vertices[v].cost = current.cost + dist
				g.graph.Vertices[v].bestVertices[0] = current.ID
				if v == dest {
					//If this is the destination update costToDest, so we can stop looking at useless Vertices
					g.costToDest = current.cost + dist
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
	for c := g.graph.Vertices[dest]; c.ID != src; c = g.graph.Vertices[c.bestVertices[0]] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return gravelmap.BestPath{Distance: g.graph.Vertices[dest].cost, Path: path}
}

// 1. Reset state
// 2. Reset the cost to destination
// set all best vertices to -1 (unused) and set the defaults *almost* as bad
func (g *dijkstraRouter) setDefaults() {
	g.destFound = false
	g.costToDest = int64(math.MaxInt64)

	for i := range g.graph.Vertices {
		g.graph.Vertices[i].bestVertices = []int{-1}
		g.graph.Vertices[i].cost = int64(math.MaxInt64) - 2
	}
}
