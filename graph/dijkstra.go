package graph

import (
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

func CalcShortest(from, to int)  {

}
