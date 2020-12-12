package path

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/thanosKontos/gravelmap"
)

func TestGmapsEncodePath(t *testing.T) {
	pts := []gravelmap.Point{
		gravelmap.Point{10.2, 11.3},
		gravelmap.Point{10.5, 11.3},
		gravelmap.Point{10.52, 11.4},
	}

	gPoly := NewGooglePolyline()
	encoded := gPoly.Encode(pts)

	assert.Equal(t, "}dg}@_`~cAary@?_|B_pR", encoded)
}

func TestGmapsEncodePathNoPoints(t *testing.T) {
	pts := []gravelmap.Point{}

	gPoly := NewGooglePolyline()
	encoded := gPoly.Encode(pts)

	assert.Equal(t, "", encoded)
}

func TestGmapsDecodePolyline(t *testing.T) {
	poly := "}dg}@_`~cAary@?_|B_pR"

	gPoly := NewGooglePolyline()
	decoded := gPoly.Decode(poly)

	assert.Len(t, decoded, 3)
}

func TestGmapsDecodeEmptyPolyline(t *testing.T) {
	poly := ""

	gPoly := NewGooglePolyline()
	decoded := gPoly.Decode(poly)

	assert.Empty(t, decoded)
}
