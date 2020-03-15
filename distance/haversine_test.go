package distance

import (
	"testing"

	"github.com/thanosKontos/gravelmap"
)

func TestHaversineDistanceBetweenPoints(t *testing.T) {
	hd := NewHaversine()

	d := hd.Distance(gravelmap.Point{Lat: 10.2, Lng: 20.2}, gravelmap.Point{Lat: 11.2, Lng: 21.2})
	if d != 155891 {
		t.Errorf("Unexpected distance %d", d)
	}

	d = hd.Distance(gravelmap.Point{Lat: 11.2, Lng: 21.2}, gravelmap.Point{Lat: 10.2, Lng: 20.2})
	if d != 155891 {
		t.Errorf("Unexpected distance %d", d)
	}

	d = hd.Distance(gravelmap.Point{Lat: 10.2001, Lng: 20.2001}, gravelmap.Point{Lat: 10.2, Lng: 20.2})
	if d != 15 {
		t.Errorf("Unexpected distance %d", d)
	}

	d = hd.Distance(gravelmap.Point{Lat: 10.2, Lng: 20.2}, gravelmap.Point{Lat: 10.2, Lng: 20.2})
	if d != 0 {
		t.Errorf("Unexpected distance %d", d)
	}
}
