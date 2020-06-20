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

	gmapsEncoder := NewGooglemaps()
	encoded := gmapsEncoder.Encode(pts)

	assert.Equal(t, "}dg}@_`~cAary@?_|B_pR", encoded)
}

func TestGmapsEncodePathNoPoints(t *testing.T) {
	pts := []gravelmap.Point{}

	gmapsEncoder := NewGooglemaps()
	encoded := gmapsEncoder.Encode(pts)

	assert.Equal(t, "", encoded)
}
