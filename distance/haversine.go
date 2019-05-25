package distance

import (
	"github.com/thanosKontos/gravelmap"
	"math"
)

var earthRadiusMetres float64 = 6371000

type Haversine struct{}

func NewHaversine() *Haversine {
	return &Haversine{}
}

func (h *Haversine) Distance(pointFrom, pointTo gravelmap.Point) (float64, error) {
	pointFrom = toRadians(pointFrom)
	pointTo = toRadians(pointTo)

	change := delta(pointFrom, pointTo)

	a := math.Pow(math.Sin(change.Lat/2), 2) + math.Cos(pointFrom.Lat)*math.Cos(pointTo.Lat)*math.Pow(math.Sin(change.Lng/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusMetres * c, nil
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func toRadians(p gravelmap.Point) gravelmap.Point {
	return gravelmap.Point{
		Lat: degreesToRadians(p.Lat),
		Lng: degreesToRadians(p.Lng),
	}
}

func delta(origin, destination gravelmap.Point) gravelmap.Point {
	return gravelmap.Point{
		Lat: origin.Lat - destination.Lat,
		Lng: origin.Lng - destination.Lng,
	}
}
