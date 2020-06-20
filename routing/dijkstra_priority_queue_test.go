package routing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/graph"
)

func TestFindShortestHappyPath(t *testing.T) {
	graph := graph.NewWeightedBidirectionalGraph()
	dijkstraRouter := NewDijsktraShortestPath(graph)

	graph.AddEdge(0, 1, 2)
	graph.AddEdge(0, 2, 1)
	graph.AddEdge(1, 3, 0)
	graph.AddEdge(2, 3, 0)
	graph.AddEdge(0, 4, 10)
	graph.AddEdge(2, 4, 5)

	bp, err := dijkstraRouter.FindShortest(0, 3)
	assert.Nil(t, err)
	assert.Equal(t, gravelmap.BestPath{Distance: 1, Path: []int{0, 2, 3}}, bp)

	bp, err = dijkstraRouter.FindShortest(0, 4)
	assert.Nil(t, err)
	assert.Equal(t, gravelmap.BestPath{Distance: 6, Path: []int{0, 2, 4}}, bp)
}

func TestFindShortestNoPathFound(t *testing.T) {
	graph := graph.NewWeightedBidirectionalGraph()
	dijkstraRouter := NewDijsktraShortestPath(graph)

	graph.AddEdge(0, 1, 1)
	graph.AddEdge(1, 2, 1)
	graph.AddEdge(2, 1, 1)
	graph.AddEdge(2, 3, 1)
	graph.AddEdge(3, 2, 1)
	graph.AddEdge(4, 3, 1)
	graph.AddEdge(4, 0, 1)

	_, err := dijkstraRouter.FindShortest(0, 4)
	assert.Equal(t, "no path found", err.Error())
}
