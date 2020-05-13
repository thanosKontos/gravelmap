package dijkstra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/graph"
)

func TestCorrectSmallGraphFindShortest(t *testing.T) {
	graph := graph.NewGraph()
	dijkstraRouter := NewDijkstra(graph)

	graph.AddVertex(0)
	graph.AddVertex(1)
	graph.AddVertex(2)
	graph.AddVertex(3)
	graph.AddVertex(4)

	graph.AddArc(0, 1, 2)
	graph.AddArc(0, 2, 1)
	graph.AddArc(1, 3, 0)
	graph.AddArc(2, 3, 0)
	graph.AddArc(0, 4, 10)
	graph.AddArc(2, 4, 5)

	bp, err := dijkstraRouter.FindShortest(0, 3)
	assert.Nil(t, err)
	assert.Equal(t, gravelmap.BestPath{Distance: 1, Path: []int{0, 2, 3}}, bp)

	bp, err = dijkstraRouter.FindShortest(0, 4)
	assert.Nil(t, err)
	assert.Equal(t, gravelmap.BestPath{Distance: 6, Path: []int{0, 2, 4}}, bp)
}

func TestCorrectLargeGraphFindShortest(t *testing.T) {
	graph := graph.NewGraph()
	dijkstraRouter := NewDijkstra(graph)

	for i := 0; i < 2000; i++ {
		v := graph.AddNewVertex()
		v.AddArc(i+1, 1)
	}
	graph.AddNewVertex()
	bp, err := dijkstraRouter.FindShortest(0, 2000)

	assert.Nil(t, err)
	assert.Equal(t, int64(2000), bp.Distance)
	assert.Equal(t, 2001, len(bp.Path))
}