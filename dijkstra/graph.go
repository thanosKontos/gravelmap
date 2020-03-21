package dijkstra

import (
	"errors"
)

//Graph contains all the graph details
type Graph struct {
	//slice of all verticies available
	Verticies []Vertex

	costToDest int64
	destFound  bool
	list       dijkstraList
}

//NewGraph creates a new empty graph
func NewGraph() *Graph {
	return &Graph{}
}

//AddNewVertex adds a new vertex at the next available index
func (g *Graph) AddNewVertex() *Vertex {
	for i, v := range g.Verticies {
		if i != v.ID {
			g.Verticies[i] = Vertex{ID: i}
			return &g.Verticies[i]
		}
	}
	return g.AddVertex(len(g.Verticies))
}

//AddVertex adds a single vertex
func (g *Graph) AddVertex(ID int) *Vertex {
	g.addVerticies(Vertex{ID: ID})
	return &g.Verticies[ID]
}

//addVerticies adds the listed verticies to the graph, overwrites any existing Vertex with the same ID.
func (g *Graph) addVerticies(verticies ...Vertex) {
	for _, v := range verticies {
		v.bestVerticies = []int{-1}
		if v.ID >= len(g.Verticies) {
			newV := make([]Vertex, v.ID+1-len(g.Verticies))
			g.Verticies = append(g.Verticies, newV...)
		}
		g.Verticies[v.ID] = v
	}
}

//AddArc is the default method for adding an arc from a Source Vertex to a Destination Vertex
func (g *Graph) AddArc(Source, Destination int, Distance int64) error {
	if len(g.Verticies) <= Source || len(g.Verticies) <= Destination {
		return errors.New("Source/Destination not found")
	}
	g.Verticies[Source].AddArc(Destination, Distance)
	return nil
}
