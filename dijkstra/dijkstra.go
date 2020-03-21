package dijkstra

import (
	"math"
)

//Shortest calculates the shortest path from src to dest
func (g *Graph) Shortest(src, dest int) (BestPath, error) {
	//Setup graph
	g.setup(src)
	return g.postSetupEvaluate(src, dest)
}

func (g *Graph) setup(src int) {
	//-1 auto list
	//Get a new list regardless
	g.forceList(-1)

	//Reset state
	g.visitedDest = false
	//Reset the best current value (worst so it gets overwritten)
	// and set the defaults *almost* as bad
	// set all best verticies to -1 (unused)

	g.setDefaults(int64(math.MaxInt64)-2, -1)
	g.best = int64(math.MaxInt64)

	//Set the distance of initial vertex 0
	g.Verticies[src].distance = 0
	//Add the source vertex to the list
	g.visiting.PushOrdered(&g.Verticies[src])
}

func (g *Graph) forceList(i int) {
	//-2 long auto
	//-1 short auto
	//0 short pq
	//1 long pq
	//2 short ll
	//3 long ll
	switch i {
	case -1:
		if len(g.Verticies) < 800 {
			g.visiting = linkedListNewLong()
			break
		} else {
			g.visiting = priorityQueueNewLong()
			break
		}
		break
	case 0:
		g.visiting = priorityQueueNewShort()
		break
	case 1:
		g.visiting = priorityQueueNewLong()
		break
	case 2:
		g.visiting = linkedListNewShort()
		break
	case 3:
		g.visiting = linkedListNewLong()
		break
	default:
		panic(i)
	}
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
	return BestPath{g.Verticies[dest].distance, path}
}

func (g *Graph) postSetupEvaluate(src, dest int) (BestPath, error) {
	var current *Vertex
	oldCurrent := -1
	for g.visiting.Len() > 0 {
		//Visit the current lowest distanced Vertex
		//TODO WTF
		current = g.visiting.PopOrdered()
		if oldCurrent == current.ID {
			continue
		}
		oldCurrent = current.ID
		//If the current distance is already worse than the best try another Vertex
		if current.distance >= g.best {
			continue
		}
		for v, dist := range current.Arcs {
			//If the arc has better access, than the current best, update the Vertex being touched
			if current.distance+dist < g.Verticies[v].distance {
				if current.bestVerticies[0] == v && g.Verticies[v].ID != dest {
					//also only do this if we aren't checkout out the best distance again
					//This seems familiar 8^)
					return BestPath{}, newErrLoop(current.ID, v)
				}
				g.Verticies[v].distance = current.distance + dist
				g.Verticies[v].bestVerticies[0] = current.ID
				if v == dest {
					//If this is the destination update best, so we can stop looking at
					// useless Verticies
					g.best = current.distance + dist
					g.visitedDest = true
					continue // Do not push if dest
				}
				//Push this updated Vertex into the list to be evaluated, pushes in
				// sorted form
				g.visiting.PushOrdered(&g.Verticies[v])
			}
		}
	}
	return g.finally(src, dest)
}

func (g *Graph) finally(src, dest int) (BestPath, error) {
	if !g.visitedDest {
		return BestPath{}, ErrNoPath
	}
	return g.bestPath(src, dest), nil
}

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

//BestPaths contains the list of best solutions
type BestPaths []BestPath
