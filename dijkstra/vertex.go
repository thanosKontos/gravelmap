package dijkstra

//Vertex is a single node in the network, contains it's ID, best cost (to
// itself from the src) and the weight to go to each other connected node (Vertex)
type Vertex struct {
	//ID of the Vertex
	ID int
	//A set of all weights to the nodes in the map
	Arcs map[int]int64

	//best cost to the Vertex
	cost          int64
	bestVerticies []int
}

//NewVertex creates a new vertex
func NewVertex(ID int) *Vertex {
	return &Vertex{ID: ID, bestVerticies: []int{-1}, Arcs: map[int]int64{}}
}

//AddArc adds an arc to the vertex, it's up to the user to make sure this is used
// correctly, firstly ensuring to use before adding to graph, or to use referenced
// of the Vertex instead of a copy. Secondly, to ensure the destination is a valid
// Vertex in the graph. Note that AddArc will overwrite any existing cost set
// if there is already an arc set to Destination.
func (v *Vertex) AddArc(Destination int, Distance int64) {
	if v.Arcs == nil {
		v.Arcs = map[int]int64{}
	}
	v.Arcs[Destination] = Distance
}
