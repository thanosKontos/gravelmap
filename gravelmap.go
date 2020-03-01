package gravelmap

// MinRoutingDistance defines the minimum distance from our start/end points to some point in our route engine
const MinRoutingDistance = 2000

const (
	WayTypePaved = iota
	WayTypeUnaved
	WayTypePath
)

type Node struct {
	GmID int
	Occurrences int
	Point Point
}

type Way struct {
	EdgeFrom int32
	EdgeTo int32
}

type WayTo struct {
	NdTo int
	WayType int8
	Grade float32
	Polyline string
}

type WayStorer interface {
	Store(ways map[int][]WayTo) error
}

type GraphWayAdder interface {
	AddWays(wayNdsOsm2GM []Node, tags map[string]string, previousLastAddedVertex int) int
}

type WayElevation struct {
	Elevations []int32
	Incline int32
	Grade float64
}

// EvaluativeWay holds info for a way to be evaluated (distance, elevation, road)
type EvaluativeWay struct {
	Tags map[string]string
	Points []Point
}

type WayCost struct {
	Cost int64
	ReverseCost int64
}

type ElevationGetterCloser interface {
	Get(points []Point, distance float64) (*WayElevation, error)
	Close()
}

type CostEvaluator interface {
	Evaluate(way EvaluativeWay) WayCost
}

type Osm2GmNodeReaderWriter interface {
	Write(osmNdID int64, gm *Node) error
	Read(osmNdID int64) *Node
}

type GMNode struct {
	ID int32
	Point
}

type GmNodeReader interface {
	Read(ndID int32) (*GMNode, error)
}

type WayPolylineReader interface {
	Read(ways []Way) []string
}

// Point represents a single point on earth
type Point struct {
	Lat float64
	Lng float64
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

// DistanceCalculator describes implementations of finding the distance between 2 points
type DistanceCalculator interface {
	Calculate(x, y Point) int64
}

type EdgeBatchStorer interface {
	BatchStore(ndBatch []GMNode) error
}

type EdgeFinder interface {
	FindClosest(point Point) (int32, error)
}

// Router describes implementations of routing between points
type Router interface {
	Route(pointFrom, pointTo Point, mode RoutingMode) ([]RoutingLeg, error)
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

// PresentableWay describes a way with all information presentable to a client
type PresentableWay struct {
	Polyline string
	SurfaceType int8
	ElevationGrade float32
}
