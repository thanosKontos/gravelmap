package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
)

func TestEmptyWeightedBidirectionalGraph(t *testing.T) {
	g := NewWeightedBidirectionalGraph()
	assert.Zero(t, len(g.Connections))
}

func TestWeightedBidirectionalGraphIndexesAllNodes(t *testing.T) {
	g := NewWeightedBidirectionalGraph()
	var ways = map[int]map[int]gravelmap.Way{
		0: {
			1: gravelmap.Way{Cost: 2},
			2: gravelmap.Way{Cost: 1},
			4: gravelmap.Way{Cost: 10},
		},
		1: {
			3: gravelmap.Way{Cost: 0},
		},
	}
	g.AddWays(ways)

	assert.Equal(t, 5, len(g.Connections))
	assert.NotEmpty(t, g.Connections[0])
	assert.NotEmpty(t, g.Connections[1])
	assert.Empty(t, g.Connections[2])
	assert.Empty(t, g.Connections[3])
	assert.Empty(t, g.Connections[4])
}
