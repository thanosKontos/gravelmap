package gravelmap

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
	EdgeTo   int32
}

type ElevationInfo struct {
	Grade float32
	From  int16
	To    int16
}

type EvaluatedWay struct {
	Points   []Point
	Tags     map[string]string
	Distance int32
	WayType  int8
	ElevationInfo
	Cost int64
}

type ElevationEvaluation struct {
	Normal  ElevationInfo
	Reverse ElevationInfo
}

type WayCost struct {
	Normal  int64
	Reverse int64
}

type WayEvaluation struct {
	Distance int32
	WayType  int8
	ElevationEvaluation
	WayCost
}

type WayAdderGetter interface {
	Add(osmNodeIds []int64, tags map[string]string)
	Get() map[int]map[int]EvaluatedWay
}

type WayStorer interface {
	Store(ways map[int]map[int]EvaluatedWay) error
}

type GraphWayAdder interface {
	AddWays(ways map[int]map[int]EvaluatedWay)
}

type WayElevation struct {
	Elevations []int32
	ElevationEvaluation
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
	Distance    int32
	Polyline    string
	SurfaceType int8
	ElevFrom    int16
	ElevTo      int16

}

type Encoder interface {
	Encode(points []Point) string
}

type PathSimplifier interface {
	Simplify(points []Point) []Point
}

type Weight struct {
	Normal  float64
	Reverse float64
}

type Weighter interface {
	WeightOffRoad(wayType int8) float64
	WeightWayAcceptance(tags map[string]string) Weight
	WeightVehicleAcceptance(tags map[string]string) float64
	WeightElevation(elevation *WayElevation) Weight
}
