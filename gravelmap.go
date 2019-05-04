package gravelmap

// MinRoutingDistance defines the minimum distance from our start/end points to some point in our routing engine
const MinRoutingDistance = 2000

type Point struct {
	Lat float64
	Lng float64
}

type Elevation interface {
	Find(lat, lng float64) (int64, error)
}

type Router interface {
	Route(pointFrom, pointTo Point) ([][]Point, error)
}
