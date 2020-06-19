package graph

import "github.com/thanosKontos/gravelmap"

// NewWeightedBidirectionalGraph creates objects of type WeightedBidirectionalGraph
func NewWeightedBidirectionalGraph() *WeightedBidirectionalGraph {
	return &WeightedBidirectionalGraph{
		Connections: make(map[int]map[int]int64),
	}
}

// WeightedBidirectionalGraph holds the connection information for a weighted bidirectional (like a road system)
type WeightedBidirectionalGraph struct {
	// Connections is a map of [fromID][toID]Weight
	Connections map[int]map[int]int64
}

func (g *WeightedBidirectionalGraph) AddEdge(fromID, toID int, weight int64) error {
	if _, ok := g.Connections[fromID]; !ok {
		g.Connections[fromID] = map[int]int64{}
	}
	if _, ok := g.Connections[toID]; !ok {
		g.Connections[toID] = map[int]int64{}
	}
	g.Connections[fromID][toID] = weight

	return nil
}

func (g *WeightedBidirectionalGraph) AddWays(ways map[int]map[int]gravelmap.Way) {
	for edgeFromId, edgeFromWays := range ways {
		for edgeToId, way := range edgeFromWays {
			g.AddEdge(edgeFromId, edgeToId, way.Cost)
		}
	}
}
