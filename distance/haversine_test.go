package distance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
)

func TestHaversineDistanceBetweenPoints(t *testing.T) {
	hd := NewHaversine()

	d := hd.Calculate(gravelmap.Point{Lat: 10.2, Lng: 20.2}, gravelmap.Point{Lat: 11.2, Lng: 21.2})
	assert.EqualValues(t, 155891, d)

	d = hd.Calculate(gravelmap.Point{Lat: 11.2, Lng: 21.2}, gravelmap.Point{Lat: 10.2, Lng: 20.2})
	assert.EqualValues(t, 155891, d)

	d = hd.Calculate(gravelmap.Point{Lat: 10.2001, Lng: 20.2001}, gravelmap.Point{Lat: 10.2, Lng: 20.2})
	assert.EqualValues(t, 15, d)

	d = hd.Calculate(gravelmap.Point{Lat: 10.2, Lng: 20.2}, gravelmap.Point{Lat: 10.2, Lng: 20.2})
	assert.Zero(t, d)
}
