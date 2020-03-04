package gravelmap

// MinRoutingDistance defines the minimum distance from our start/end points to some point in our route engine
const MinRoutingDistance = 2000

const (
	WayTypePaved int8 = iota
	WayTypeUnaved
	WayTypePath
)

type Node struct {
	Id          int
	Occurrences int
	Point       Point
}

type Way struct {
	EdgeFrom int32
	EdgeTo int32
}

type ElevationInfo struct {
	Grade float32
	From int16
	To int16
}

type WayTo struct {
	NdTo int
	Points []Point
	Tags map[string]string
	Distance int32
	WayType int8
	ElevationInfo
	Cost int64
}

type ElevationEvaluation struct {
	Normal ElevationInfo
	Reverse ElevationInfo
}

type WayCost struct {
	Normal int64
	Reverse int64
}

type WayEvaluation struct {
	Distance int32
	WayType int8
	ElevationEvaluation
	WayCost
}

type WayAdderGetter interface {
	Add(osmNodeIds []int64, tags map[string]string)
	Get() map[int][]WayTo
}

type WayStorer interface {
	Store(ways map[int][]WayTo) error
}

type GraphWayAdder interface {
	AddWays(ways map[int][]WayTo)
}

type WayElevation struct {
	Elevations []int32
	ElevationEvaluation
}

// EvaluativeWay holds info for a way to be evaluated (distance, elevation, road)
type EvaluativeWay struct {
	Tags map[string]string
	Points []Point
}

type ElevationGetterCloser interface {
	Get(points []Point, distance float64) (*WayElevation, error)
	Close()
}

type CostEvaluator interface {
	Evaluate(points []Point, tags map[string]string) WayEvaluation
}

type Osm2GmNodeReaderWriter interface {
	Write(osmNdID int64, gm *Node) error
	Read(osmNdID int64) *Node
}

type GmNodeReader interface {
	Read(ndID int) (*Node, error)
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
	BatchStore(ndBatch []Node) error
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
	Distance int32
	Polyline string
	SurfaceType int8
	ElevationInfo
}

type Encoder interface {
	Encode(points []Point) string
}
