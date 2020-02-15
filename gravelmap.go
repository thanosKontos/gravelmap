package gravelmap

// MinRoutingDistance defines the minimum distance from our start/end points to some point in our route engine
const MinRoutingDistance = 2000

// Point represents a single point on earth
type Point struct {
	Lat float64
	Lng float64
}

// WayElevation represents the elevation of a road (start, end and the gradient percentage)
type WayElevation struct {
	Grade  float64
	Start  float64
	End    float64
	Length float64
}

// RoutingLegElevation represents the elevation routing leg
type RoutingLegElevation struct {
	Grade float64
	Start float64
	End   float64
}

// RoutingLeg represents individual parts of a route
type RoutingLeg struct {
	Coordinates []Point
	Length      float64
	Paved       bool
	Elevation   *RoutingLegElevation
}

// RoutingMode is a lookup value of the different routing modes
type RoutingMode int

const (
	Normal RoutingMode = iota
	OnlyUnpavedAccountElevation
	OnlyUnpavedHardcore
	NoLengthCareEasiest
	NoLengthCareNormal
	NoLengthOnlyUnpavedHardcore
)

// Importer describes implementations of import raw routing data
type Importer interface {
	Import() error
}

// ElevationFinder describes implementations of finding the elevation for a single point
type ElevationFinder interface {
	FindElevation(Point) (float64, error)
}

// DistanceFinder describes implementations of finding the distance between 2 points
type DistanceFinder interface {
	Distance(pointFrom, pointTo Point) (float64, error)
}

// ElevationGrader describes implementations of finding the elevation for continuous points
type ElevationGrader interface {
	Grade([]Point, float64) (*WayElevation, error)
}

// WayGrader describes implementations of grading the elevation of roads/paths
type WayGrader interface {
	GradeWays() error
}

// Router describes implementations of routing between points
type Router interface {
	Route(pointFrom, pointTo Point, mode RoutingMode) ([]RoutingLeg, error)
	Close() error
}

// RouterPreparer describes implementations of preparing the routing (creating graphs, files etc)
type RouterPreparer interface {
	Prepare() error
	Close() error
}

// OsmFilter describes implementations of filtering the useless OSM data
type OsmFilter interface {
	Filter() error
}

// Logger describes implementations of logging
type Logger interface {
	Info(log interface{})
	Debug(log interface{})
	Warning(log interface{})
	Error(log interface{})
}

type OsmNode struct {
	NdID int64
}

type OsmWay struct {
	WayId int64
	NdIds []int64
}

type OsmIterator interface {
	Iterate() (chan OsmNode, chan OsmWay)
}
