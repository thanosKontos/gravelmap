package dijkstra

import (
	"math"
)

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

//Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	g.setup(src)
	return g.postSetupEvaluate(src, dest)
}

func (g *Graph) setup(src int) {
	g.setupList()
	g.setDefaults()

	//Set the cost of initial vertex 0 and add it to the list
	g.Verticies[src].cost = 0
	g.list.PushOrdered(&g.Verticies[src])
}

func (g *Graph) setupList() {
	if len(g.Verticies) < 800 {
		g.list = linkedListNewLong()
		return
	}

	g.list = priorityQueueNewLong()
	return
}

func (g *Graph) postSetupEvaluate(src, dest int) (BestPath, error) {
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
			if current.cost+dist < g.Verticies[v].cost {
				if current.bestVerticies[0] == v && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best cost again
					//This seems familiar 8^)
					return BestPath{}, newErrLoop(current.ID, v)
				}
				g.Verticies[v].cost = current.cost + dist
				g.Verticies[v].bestVerticies[0] = current.ID
				if v == dest {
					//If this is the destination update costToDest, so we can stop looking at useless Verticies
					g.costToDest = current.cost + dist
					g.destFound = true
					continue // Do not push if dest
				}

				//Push this updated Vertex into the list to be evaluated, pushes in sorted form
				g.list.PushOrdered(&g.Verticies[v])
			}
		}
	}
	return g.finally(src, dest)
}

func (g *Graph) finally(src, dest int) (BestPath, error) {
	if !g.destFound {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}

func (g *Graph) bestPath(src, dest int) BestPath {
	var path []int
	for c := g.Verticies[dest]; c.ID != src; c = g.Verticies[c.bestVerticies[0]] {
		path = append(path, c.ID)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return BestPath{g.Verticies[dest].cost, path}
}

// 1. Reset state
// 2. Reset the cost to destination
// set all best verticies to -1 (unused) and set the defaults *almost* as bad
func (g *Graph) setDefaults() {
	g.destFound = false
	g.costToDest = int64(math.MaxInt64)

	for i := range g.Verticies {
		g.Verticies[i].bestVerticies = []int{-1}
		g.Verticies[i].cost = int64(math.MaxInt64) - 2
	}
}
