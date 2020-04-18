package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewVertex(t *testing.T) {
	v := NewVertex(10)

	assert.Equal(t, 10, v.ID)
	assert.NotNil(t, v.Arcs)
}
