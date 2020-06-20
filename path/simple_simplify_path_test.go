package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/distance"
)

func TestSimplePathSimplification(t *testing.T) {
	points := []gravelmap.Point{
		{38.06988, 22.98194},
		{38.06994, 22.98195},
		{38.06999, 22.98197},
		{38.07005, 22.98201},
		{38.07009, 22.98206},
		{38.07014, 22.98209},
		{38.07026, 22.98215},
		{38.07038, 22.98218},
		{38.07045, 22.98223},
		{38.07053, 22.98235},
		{38.07067, 22.98253},
		{38.07075, 22.98268},
		{38.07084, 22.98285},
		{38.0709, 22.983},
		{38.07094, 22.98312},
		{38.07101, 22.98321},
		{38.07111, 22.98332},
		{38.07127, 22.9835},
		{38.07142, 22.98369},
		{38.07158, 22.98386},
		{38.07161, 22.98389},
	}

	hd := distance.NewHaversine()
	dp := NewSimpleSimplifiedPath(hd)
	simplified := dp.Simplify(points)

	assert.Equal(t, 11, len(simplified))
}

func TestSimplePathSimplificationAnother(t *testing.T) {
	points := []gravelmap.Point{
		{37.9818586, 23.8148583},
		{37.9818, 23.814885},
		{37.981713, 23.814965},
		{37.981556, 23.815165},
		{37.981465, 23.815331},
		{37.981323, 23.815667},
		{37.981242, 23.815939},
		{37.981208, 23.816112},
		{37.981207, 23.8162},
		{37.981245, 23.816293},
		{37.981296, 23.816334},
		{37.981374, 23.816343},
		{37.981509, 23.816291},
		{37.981882, 23.816092},
		{37.98195, 23.816079},
		{37.981994, 23.816094},
		{37.982023, 23.816144},
		{37.982005, 23.81623},
		{37.981951, 23.816291},
		{37.981677, 23.816489},
		{37.981198, 23.81701},
		{37.981096, 23.817099},
		{37.981038, 23.81712},
		{37.980952, 23.817111},
		{37.980898, 23.817074},
		{37.980586, 23.816611},
		{37.98048, 23.816492},
		{37.980426, 23.816457},
		{37.980317, 23.816451},
		{37.980239, 23.8165},
		{37.980044, 23.816714},
		{37.9799365, 23.8168567},
		{37.979927, 23.8169186},
		{37.9799506, 23.8169804},
	}

	hd := distance.NewHaversine()
	dp := NewSimpleSimplifiedPath(hd)
	simplified := dp.Simplify(points)

	assert.Equal(t, 30, len(simplified))
}