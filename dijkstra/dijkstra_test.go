package dijkstra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectSmallGraphFindShortest(t *testing.T) {
	graph := NewGraph()

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

	bp, err := graph.Shortest(0, 3)
	assert.Nil(t, err)
	assert.Equal(t, BestPath{1, []int{0, 2, 3}}, bp)

	bp, err = graph.Shortest(0, 4)
	assert.Nil(t, err)
	assert.Equal(t, BestPath{6, []int{0, 2, 4}}, bp)
}

func TestCorrectLargeGraphFindShortest(t *testing.T) {
	g := NewGraph()
	for i := 0; i < 2000; i++ {
		v := g.AddNewVertex()
		v.AddArc(i+1, 1)
	}
	g.AddNewVertex()
	bp, err := g.Shortest(0, 2000)

	assert.Nil(t, err)
	assert.Equal(t, int64(2000), bp.Distance)
	assert.Equal(t, 2001, len(bp.Path))
}
