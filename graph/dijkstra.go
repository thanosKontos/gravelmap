package graph

import (
	"github.com/thanosKontos/gravelmap"
	dijkstra2 "github.com/thanosKontos/gravelmap/dijkstra"
)

type dijkstra struct {
	costEvaluator gravelmap.CostEvaluator
	graph *dijkstra2.Graph
}

func NewDijkstra(costEvaluator gravelmap.CostEvaluator) *dijkstra {
	return &dijkstra{
		costEvaluator: costEvaluator,
		graph: dijkstra2.NewGraph(),
	}
}

func (d *dijkstra) Get() *dijkstra2.Graph {
	return d.graph
}

func (d *dijkstra) AddWays(wayNdsOsm2GM []gravelmap.NodeOsm2GM, tags map[string]string, previousLastAddedVertex int) int {
	var evaluativeWay = gravelmap.EvaluativeWay{Tags: tags}
	var previousSubwayPoint = gravelmap.Point{}
	var previousEdge gravelmap.NodeOsm2GM
	var firstEdge gravelmap.NodeOsm2GM
	lastAddedVertex := -1

	for _, ndOsm2GM := range wayNdsOsm2GM {
		if isEdge := ndOsm2GM.Occurrences > 1; isEdge {
			if ndOsm2GM.GmID > previousLastAddedVertex || previousLastAddedVertex == 0 {
				d.graph.AddVertex(ndOsm2GM.GmID)
				lastAddedVertex = ndOsm2GM.GmID
			}

			if isFirstEdge := firstEdge == (gravelmap.NodeOsm2GM{}); isFirstEdge {
				evaluativeWay.Points = append(evaluativeWay.Points, ndOsm2GM.Point)
				previousSubwayPoint = ndOsm2GM.Point
				firstEdge = ndOsm2GM
				previousEdge = ndOsm2GM
				continue
			}

			evaluativeWay.Points = append(evaluativeWay.Points, ndOsm2GM.Point)

			cost := d.costEvaluator.Evaluate(evaluativeWay)

			d.graph.AddArc(ndOsm2GM.GmID, previousEdge.GmID, cost.Cost)
			d.graph.AddArc(previousEdge.GmID, ndOsm2GM.GmID, cost.ReverseCost)

			evaluativeWay.Points = []gravelmap.Point{ndOsm2GM.Point}

			previousEdge = ndOsm2GM
			previousSubwayPoint = ndOsm2GM.Point
		} else {
			if hasPreviousSubwayPoint := previousSubwayPoint != (gravelmap.Point{}); hasPreviousSubwayPoint {
				evaluativeWay.Points = append(evaluativeWay.Points, ndOsm2GM.Point)
			}
		}
	}

	return lastAddedVertex
}