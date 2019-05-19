package gravelmap

// MinRoutingDistance defines the minimum distance from our start/end points to some point in our routing engine
const MinRoutingDistance = 2000

type Point struct {
	Lat float64
	Lng float64
}

type Elevation interface {
	Import() error
}

type Router interface {
	Route(pointFrom, pointTo Point) ([][]Point, error)
	Close() error
}

type RouterPreparer interface {
	Prepare() error
	Close() error
}

type RouterImporter interface {
	Import() error
}

type OsmFilter interface {
	Filter() error
}

type Logger interface {
	Info(log interface{})
	Debug(log interface{})
	Warning(log interface{})
	Error(log interface{})
}
