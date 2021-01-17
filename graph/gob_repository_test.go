package graph

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
)

func TestFetchAndRetrieveGraph(t *testing.T) {
	defer os.Remove("../fixtures/graph_abc.gob")
	repo := NewGobRepo("../fixtures")

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

	err := repo.Store(g, "abc")
	assert.Nil(t, err)

	fetchedGraph, err := repo.Fetch("abc")
	assert.Nil(t, err)
	assert.Equal(t, 5, len(fetchedGraph.Connections))
	assert.NotEmpty(t, fetchedGraph.Connections[0])
	assert.NotEmpty(t, fetchedGraph.Connections[1])
	assert.Empty(t, fetchedGraph.Connections[2])
	assert.Empty(t, fetchedGraph.Connections[3])
	assert.Empty(t, fetchedGraph.Connections[4])
}

func TestErrorRetrieveNonExistingGraph(t *testing.T) {
	repo := NewGobRepo("../fixtures")
	_, err := repo.Fetch("xyz")
	assert.NotNil(t, err)
}
