package distance

import (
	"math"

	"github.com/thanosKontos/gravelmap"
)

const (
	earthRadius = 6371000
)

type haversine struct{}

func NewHaversine() *haversine {
	return &haversine{}
}

func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// Distance calculates the distance between two points on earth on a straight line (in meters)
func (h *haversine) Calculate(x, y gravelmap.Point) int64 {
	lat1 := degreesToRadians(x.Lat)
	lng1 := degreesToRadians(x.Lng)
	lat2 := degreesToRadians(y.Lat)
	lng2 := degreesToRadians(y.Lng)

	diffLat := lat2 - lat1
	diffLon := lng2 - lng1

	a := math.Pow(math.Sin(diffLat/2), 2) + math.Cos(lat1)*math.Cos(lat2) * math.Pow(math.Sin(diffLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int64(c * earthRadius)
}
