package gravelmap

const (
	WayTypePaved int8 = iota
	WayTypeUnpaved
	WayTypePath
)

// NodePoint is the data that holds the coordinates of a node
type NodePoint struct {
	NodeID int32
	Point  Point
}

type ConnectionNode struct {
	ID            int
	ConnectionCnt int
	Point         Point
}

type Edge struct {
	NodeFrom int32
	NodeTo   int32
}

type ElevationInfo struct {
	Grade float32
	From  int16
	To    int16
}

type Way struct {
	Points   []Point
	Tags     map[string]string
	Distance int32
	Type     int8
	ElevationInfo
	Cost int64

	// for debug reasons not needed really for production code
	OriginalOsmID int64
}

type Costable interface {
	GetCost() int64
}

func (w *Way) GetCost() int64 {
	return w.Cost
}

type BidirectionalElevationInfo struct {
	Normal  ElevationInfo
	Reverse ElevationInfo
}

type BidirectionalCost struct {
	Normal  int64
	Reverse int64
}

type WayEvaluation struct {
	Distance int32
	WayType  int8
	BidirectionalElevationInfo
	BidirectionalCost
}

type WayAdderGetter interface {
	Add(osmNodeIds []int64, tags map[string]string, osmID int64)
	Get() map[int]map[int]Way
}

type WayStorer interface {
	Store(ways map[int]map[int]Way) error
}

type GraphWayAdder interface {
	AddWays(ways map[int]map[int]Way)
}

type WayElevation struct {
	Elevations []int32
	BidirectionalElevationInfo
}

type ElevationGetterCloser interface {
	Get(points []Point, distance float64) (*WayElevation, error)
	Close()
}

type CostEvaluator interface {
	Evaluate(points []Point, tags map[string]string) WayEvaluation
}

type Osm2GmNodeReaderWriter interface {
	Write(osmNdID int64, gm *ConnectionNode) error
	Read(osmNdID int64) *ConnectionNode
}

type Osm2LatLngWriter interface {
	Write(osmID int, point Point)
}

type Osm2LatLngReader interface {
	Read(ndID int) (Point, error)
}

// Point represents a single point on earth
type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// ElevationRange represents the elevation routing leg
type ElevationRange struct {
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// RoutingLeg represents individual parts of a route
type RoutingLeg struct {
	Coordinates []Point         `json:"points"`
	Length      float64         `json:"distance"`
	WayType     string          `json:"type"`
	Elevation   *ElevationRange `json:"elev"`
	OsmID       int64           `json:"osm_id"`
}

// DistanceCalculator describes implementations of finding the distance between 2 points
type DistanceCalculator interface {
	Calculate(x, y Point) int64
}

// NodePointStorer describe implementations of storing NodePoints
type NodePointStorer interface {
	BatchStore(nds []NodePoint) error
}

// EdgeFinder describe implementations of finding nodes from coordinates
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
	OsmID       int64
}

type Encoder interface {
	Encode(points []Point) string
}

type Decoder interface {
	Decode(string) []Point
}

type PathSimplifier interface {
	Simplify(points []Point) []Point
}

type BidirectionalWeight struct {
	Normal  float64
	Reverse float64
}

type Weighter interface {
	WeightOffRoad(wayType int8) float64
	WeightWayAcceptance(tags map[string]string) BidirectionalWeight
	WeightVehicleAcceptance(tags map[string]string) float64
	WeightElevation(elevation *WayElevation) BidirectionalWeight
}

//BestPath contains the solution of the most optimal path
type BestPath struct {
	Distance int64
	Path     []int
}

type ShortestFinder interface {
	FindShortest(src, dest int) (BestPath, error)
}

type EdgeReader interface {
	Read(edges []Edge) ([]PresentableWay, error)
}
