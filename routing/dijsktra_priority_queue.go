package routing

import (
	"math"

	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/graph"
)

// TODO: in the future would be nice to connect dijkstra with a graph interface instead of concrete implementation
// as in here. Not going to need abstraction soon, so it is ok for now (integration tests are also ok at the moment)
type dijsktraShortestPath struct {
	graph             *graph.WeightedBidirectionalGraph
	list              DijkstraList
	destFound         bool
	costToDest        int64
	evaluatedVertices map[int]evaluatedVertex
}

type evaluatedVertex struct {
	id           int
	cost         int64
	bestVertices []int
}

func NewDijsktraShortestPath(g *graph.WeightedBidirectionalGraph) dijsktraShortestPath {
	return dijsktraShortestPath{
		graph: g,
	}
}

func (dsp dijsktraShortestPath) FindShortest(fromID, toID int) (gravelmap.BestPath, error) {
	//setup
	dsp.destFound = false
	dsp.costToDest = int64(math.MaxInt64)
	dsp.list = PriorityQueueNewLong()
	dsp.evaluatedVertices = make(map[int]evaluatedVertex)
	for i := range dsp.graph.Connections {
		dsp.evaluatedVertices[i] = evaluatedVertex{
			id:           i,
			cost:         int64(math.MaxInt64),
			bestVertices: []int{-1},
		}
	}
	fromVtx := dsp.evaluatedVertices[fromID]
	fromVtx.cost = 0
	dsp.list.PushOrdered(&fromVtx)

	var current *evaluatedVertex
	oldCurrent := -1
	for dsp.list.Len() > 0 {
		current = dsp.list.PopOrdered()
		if oldCurrent == current.id {
			continue
		}
		oldCurrent = current.id
		//If the current cost is already worse than the best one try another node
		if current.cost >= dsp.costToDest {
			continue
		}

		currNodeConns := dsp.graph.Connections[current.id]
		for v, dist := range currNodeConns {
			//If the arc has better access, than the current costToDest, update the Vertex being touched
			if current.cost+dist < dsp.evaluatedVertices[v].cost {
				if current.bestVertices[0] == v && dsp.evaluatedVertices[v].id != toID {
					//also only do this if we aren't checkout out the best cost again
					//This seems familiar 8^)
					return gravelmap.BestPath{}, newErrLoop(current.id, v)
				}

				bv := dsp.evaluatedVertices[v].bestVertices
				bv[0] = current.id
				dsp.evaluatedVertices[v] = evaluatedVertex{id: v, cost: current.cost + dist, bestVertices: bv}
				if v == toID {
					//If this is the destination update costToDest, so we can stop looking at useless Vertices
					dsp.costToDest = current.cost + dist
					dsp.destFound = true
					continue // Do not push if dest
				}

				//Push this updated Vertex into the list to be evaluated, pushes in sorted form
				toBeEvaluatedVtx := dsp.evaluatedVertices[v]
				toBeEvaluatedVtx.bestVertices = []int{-1}
				dsp.list.PushOrdered(&toBeEvaluatedVtx)
			}
		}
	}

	if dsp.destFound {
		return dsp.getCalculatedBestPath(fromID, toID), nil
	}

	return gravelmap.BestPath{}, ErrNoPath
}

func (dsp *dijsktraShortestPath) getCalculatedBestPath(src, dest int) gravelmap.BestPath {
	var path []int
	for c := dsp.evaluatedVertices[dest]; c.id != src; c = dsp.evaluatedVertices[c.bestVertices[0]] {
		path = append(path, c.id)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return gravelmap.BestPath{Distance: dsp.evaluatedVertices[dest].cost, Path: path}
}
