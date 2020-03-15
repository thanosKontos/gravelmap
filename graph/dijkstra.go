package graph

import (
	"github.com/thanosKontos/gravelmap"
	dijkstra2 "github.com/thanosKontos/gravelmap/dijkstra"
)

type dijkstra struct {
	graph *dijkstra2.Graph
}

func NewDijkstra() *dijkstra {
	return &dijkstra{
		graph: dijkstra2.NewGraph(),
	}
}

func (d *dijkstra) Get() *dijkstra2.Graph {
	return d.graph
}

var alreadyAddedNodes = map[int]struct{}{}

func (d *dijkstra) AddWays(ways map[int]map[int]gravelmap.EvaluatedWay) {
	for edgeFromId, edgeFromWays := range ways {
		for edgeToId, way := range edgeFromWays {
			if _, ok := alreadyAddedNodes[edgeFromId]; !ok {
				d.graph.AddVertex(edgeFromId)
			}

			if _, ok := alreadyAddedNodes[edgeToId]; !ok {
				d.graph.AddVertex(edgeToId)
			}

			alreadyAddedNodes[edgeFromId] = struct{}{}
			alreadyAddedNodes[edgeToId] = struct{}{}

			d.graph.AddArc(edgeFromId, edgeToId, way.Cost)
		}
	}
}
