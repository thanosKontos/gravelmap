package graph

import (
	"errors"
)

//Graph contains all the graph details
type Graph struct {
	//slice of all vertices available
	Vertices []Vertex
}

//NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{}
}

//AddNewVertex adds a new vertex at the next available index
func (g *Graph) AddNewVertex() *Vertex {
	for i, v := range g.Vertices {
		if i != v.ID {
			g.Vertices[i] = Vertex{ID: i}
			return &g.Vertices[i]
		}
	}
	return g.AddVertex(len(g.Vertices))
}

//AddVertex adds a single vertex
func (g *Graph) AddVertex(ID int) *Vertex {
	v := Vertex{ID: ID}
	v.BestVertices = []int{-1}
	if v.ID >= len(g.Vertices) {
		newV := make([]Vertex, v.ID+1-len(g.Vertices))
		g.Vertices = append(g.Vertices, newV...)
	}
	g.Vertices[v.ID] = v

	return &g.Vertices[ID]
}

//AddArc is the default method for adding an arc from a Source Vertex to a Destination Vertex
func (g *Graph) AddArc(Source, Destination int, Distance int64) error {
	if len(g.Vertices) <= Source || len(g.Vertices) <= Destination {
		return errors.New("Source/Destination not found")
	}
	g.Vertices[Source].AddArc(Destination, Distance)
	return nil
}
