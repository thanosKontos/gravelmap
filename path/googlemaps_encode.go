package path

import (
	"github.com/thanosKontos/gravelmap"
	"googlemaps.github.io/maps"
)

type googlemaps struct {
}

func NewGooglemaps() *googlemaps {
	return &googlemaps{}
}

func (gm *googlemaps) Encode(points []gravelmap.Point) string {
	var latLngs []maps.LatLng
	for _, pt := range points {
		latLngs = append(latLngs, maps.LatLng{Lat: pt.Lat, Lng: pt.Lng})
	}

	return maps.Encode(latLngs)
}
