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

func (d *dijkstra) AddWays(ways map[int][]gravelmap.WayTo) {
	for edgeFromId, waysTo := range ways {
		for _, way := range waysTo {
			if _, ok := alreadyAddedNodes[edgeFromId]; !ok {
				d.graph.AddVertex(edgeFromId)
			}

			if _, ok := alreadyAddedNodes[way.NdTo]; !ok {
				d.graph.AddVertex(way.NdTo)
			}

			alreadyAddedNodes[edgeFromId] = struct{}{}
			alreadyAddedNodes[way.NdTo] = struct{}{}

			d.graph.AddArc(edgeFromId, way.NdTo, way.Cost)
		}
	}
}
