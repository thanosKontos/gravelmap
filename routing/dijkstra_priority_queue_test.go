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

	var ways = map[int]map[int]gravelmap.Way{
		0: {
			1: gravelmap.Way{Cost: 2},
			2: gravelmap.Way{Cost: 1},
			4: gravelmap.Way{Cost: 10},
		},
		1: {
			3: gravelmap.Way{Cost: 0},
		},
		2: {
			3: gravelmap.Way{Cost: 0},
			4: gravelmap.Way{Cost: 5},
		},
	}
	graph.AddWays(ways)

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

	var ways = map[int]map[int]gravelmap.Way{
		0: {
			1: gravelmap.Way{Cost: 1},
		},
		1: {
			2: gravelmap.Way{Cost: 1},
		},
		2: {
			1: gravelmap.Way{Cost: 1},
			3: gravelmap.Way{Cost: 1},
		},
		3: {
			2: gravelmap.Way{Cost: 1},
		},
		4: {
			3: gravelmap.Way{Cost: 1},
			0: gravelmap.Way{Cost: 1},
		},
	}
	graph.AddWays(ways)

	_, err := dijkstraRouter.FindShortest(0, 4)
	assert.Equal(t, "no path found", err.Error())
}
