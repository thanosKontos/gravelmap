package dijkstra

import "github.com/thanosKontos/gravelmap"

var alreadyAddedNodes = map[int]struct{}{}

func (g *Graph) AddWays(ways map[int]map[int]gravelmap.Way) {
	for edgeFromId, edgeFromWays := range ways {
		for edgeToId, way := range edgeFromWays {
			if _, ok := alreadyAddedNodes[edgeFromId]; !ok {
				g.AddVertex(edgeFromId)
			}

			if _, ok := alreadyAddedNodes[edgeToId]; !ok {
				g.AddVertex(edgeToId)
			}

			alreadyAddedNodes[edgeFromId] = struct{}{}
			alreadyAddedNodes[edgeToId] = struct{}{}

			g.AddArc(edgeFromId, edgeToId, way.Cost)
		}
	}
}
