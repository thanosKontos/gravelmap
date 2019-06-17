package srtm_ascii

import (
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/mock"
	"testing"
)

func TestElevationOfTwoPoints(t *testing.T) {
	eg, _ := NewElevationGrader(mock.NewElevationMock())
	points := []gravelmap.Point{
		{Lat: 10.2, Lng: 20.2},
		{Lat: 11.2, Lng: 21.2},
	}
	wayEle, _ := eg.Grade(points, 50)
	expectedWayEle := gravelmap.WayElevation{0, 10, 10, 50}

	if expectedWayEle != *wayEle {
		t.Errorf("Not expected way elevation %v instead of %v", wayEle, expectedWayEle)
	}
}
