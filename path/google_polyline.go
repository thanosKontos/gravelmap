package path

import (
	"github.com/thanosKontos/gravelmap"
	"googlemaps.github.io/maps"
)

type googlePolyline struct {
}

func NewGooglePolyline() *googlePolyline {
	return &googlePolyline{}
}

func (gm *googlePolyline) Encode(points []gravelmap.Point) string {
	var latLngs []maps.LatLng
	for _, pt := range points {
		latLngs = append(latLngs, maps.LatLng{Lat: pt.Lat, Lng: pt.Lng})
	}

	return maps.Encode(latLngs)
}

func (gm *googlePolyline) Decode(polyline string) []gravelmap.Point {
	var latLngs []gravelmap.Point
	tmpLatLngs, _ := maps.DecodePolyline(polyline)

	for _, latlng := range tmpLatLngs {
		latLngs = append(latLngs, gravelmap.Point{Lat: latlng.Lat, Lng: latlng.Lng})
	}

	return latLngs
}
