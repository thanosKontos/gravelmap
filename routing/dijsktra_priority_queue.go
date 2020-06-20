package routing

import (
	"math"

	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/graph"
)

// TODO: in the future would be nice to connect dijkstra with a graph interface instead of concrete implementation
// as in here. Not going to need abstraction soon, so it is ok for now (integration tests are also ok at the moment)
type dijsktraShortestPath struct {
	graph          *graph.WeightedBidirectionalGraph
	list           DijkstraList
	destFound      bool
	costToDest     int64
	evaluatedNodes map[int]evaluatedNode
}

type evaluatedNode struct {
	id        int
	cost      int64
	bestNodes []int
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
	dsp.evaluatedNodes = make(map[int]evaluatedNode)
	for i := range dsp.graph.Connections {
		dsp.evaluatedNodes[i] = evaluatedNode{
			id:        i,
			cost:      int64(math.MaxInt64),
			bestNodes: []int{-1},
		}
	}
	fromVtx := dsp.evaluatedNodes[fromID]
	fromVtx.cost = 0
	dsp.list.PushOrdered(&fromVtx)

	var current *evaluatedNode
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
			//If the arc has better access, than the current costToDest, update the node being touched
			if current.cost+dist < dsp.evaluatedNodes[v].cost {
				if current.bestNodes[0] == v && dsp.evaluatedNodes[v].id != toID {
					//also only do this if we aren't checkout out the best cost again
					//This seems familiar 8^)
					return gravelmap.BestPath{}, newErrLoop(current.id, v)
				}

				bv := dsp.evaluatedNodes[v].bestNodes
				bv[0] = current.id
				dsp.evaluatedNodes[v] = evaluatedNode{id: v, cost: current.cost + dist, bestNodes: bv}
				if v == toID {
					//If this is the destination update costToDest, so we can stop looking at useless nodes
					dsp.costToDest = current.cost + dist
					dsp.destFound = true
					continue // Do not push if dest
				}

				//Push this updated node into the list to be evaluated
				toBeEvaluatedVtx := dsp.evaluatedNodes[v]
				toBeEvaluatedVtx.bestNodes = []int{-1}
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
	for c := dsp.evaluatedNodes[dest]; c.id != src; c = dsp.evaluatedNodes[c.bestNodes[0]] {
		path = append(path, c.id)
	}
	path = append(path, src)
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return gravelmap.BestPath{Distance: dsp.evaluatedNodes[dest].cost, Path: path}
}
