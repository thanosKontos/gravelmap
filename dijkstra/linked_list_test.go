package dijkstra

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLazyInit(t *testing.T) {
	ll := linkedList{root: element{}, len: 0}

	assert.Nil(t, ll.root.next)
	assert.Nil(t, ll.root.prev)

	ll.lazyinit()
	assert.NotNil(t, ll.root.next)
	assert.NotNil(t, ll.root.prev)
}

func TestEmptyList(t *testing.T) {
	ll := linkedList{root: element{}, len: 0}

	assert.Nil(t, ll.front())
	assert.Nil(t, ll.back())
}
