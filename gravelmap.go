package gravelmap

// MinRoutingDistance defines the minimum distance from our start/end points to some point in our route engine
const MinRoutingDistance = 2000

type Point struct {
	Lat float64
	Lng float64
}

type RoutingFeature struct {
	Type string
	Coordinates []Point
	Options struct{
		OSMID int64
		ElevationCost float64
	}
}

type Importer interface {
	Import() error
}

type ElevationFinder interface {
	FindElevation(Point) (float64, error)
}

type DistanceFinder interface {
	Distance(pointFrom, pointTo Point) (float64, error)
}

type ElevationGrader interface {
	Grade([]Point) (float64, error)
}

type Router interface {
	Route(pointFrom, pointTo Point) ([]RoutingFeature, error)
	Close() error
}

type RouterPreparer interface {
	Prepare() error
	Close() error
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
