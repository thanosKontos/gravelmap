package dijkstra

import (
	"errors"
	"fmt"
)

//Graph contains all the graph details
type Graph struct {
	best        int64
	visitedDest bool
	//slice of all verticies available
	Verticies       []Vertex
	list            dijkstraList
	highestMapIndex int
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
	g.AddVerticies(Vertex{ID: ID})
	return &g.Verticies[ID]
}

//GetVertex gets the reference of the specified vertex. An error is thrown if
// there is no vertex with that index/ID.
func (g *Graph) GetVertex(ID int) (*Vertex, error) {
	if ID >= len(g.Verticies) {
		return nil, errors.New("vertex not found")
	}
	return &g.Verticies[ID], nil
}

//AddArc is the default method for adding an arc from a Source Vertex to a Destination Vertex
func (g *Graph) AddArc(Source, Destination int, Distance int64) error {
	if len(g.Verticies) <= Source || len(g.Verticies) <= Destination {
		return errors.New("Source/Destination not found")
	}
	g.Verticies[Source].AddArc(Destination, Distance)
	return nil
}

func (g Graph) validate() error {
	for _, v := range g.Verticies {
		for a := range v.Arcs {
			if a >= len(g.Verticies) || (g.Verticies[a].ID == 0 && a != 0) {
				return errors.New(fmt.Sprint("Graph validation error;", "Vertex ", a, " referenced in arcs by Vertex ", v.ID))
			}
		}
	}
	return nil
}

//SetDefaults sets the distance and best node to that specified
func (g *Graph) setDefaults(Distance int64, BestNode int) {
	for i := range g.Verticies {
		g.Verticies[i].bestVerticies = []int{BestNode}
		g.Verticies[i].distance = Distance
	}
}
