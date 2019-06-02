package gravelmap

// MinRoutingDistance defines the minimum distance from our start/end points to some point in our route engine
const MinRoutingDistance = 2000

type Point struct {
	Lat float64
	Lng float64
}

type WayElevation struct {
	Grade  float64
	Start  float64
	End    float64
	Length float64
}

type RoutingLegElevation struct {
	Grade float64
	Start float64
	End   float64
}

type RoutingLeg struct {
	Coordinates []Point
	Length      float64
	Paved       bool
	Elevation   *RoutingLegElevation
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
	Grade([]Point, float64) (*WayElevation, error)
}

type Router interface {
	Route(pointFrom, pointTo Point) ([]RoutingLeg, error)
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
