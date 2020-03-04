package dijkstra

import "testing"

func TestNewVertex(t *testing.T) {
	v := NewVertex(10)
	if v.ID != 10 {
		t.Error("NewVertex ID not set")
	}
	if v.Arcs == nil {
		t.Error("NewVertex arcs map not initialised")
	}
}
