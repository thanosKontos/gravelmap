package hgt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
)

func TestGetElevation(t *testing.T) {
	f, err := os.Open("../../fixtures/N37E023.hgt")
	assert.Nil(t, err)

	srtm := NewStrm1(f)

	ele, err := srtm.Get(gravelmap.Point{Lat: 37.9478199, Lng: 23.815421})
	assert.Nil(t, err)
	assert.Equal(t, int32(1010), ele) // Mt imittos summit

	ele, err = srtm.Get(gravelmap.Point{Lat: 37.7310614, Lng: 23.9364142})
	assert.Nil(t, err)
	assert.Equal(t, int32(5), ele) // Seaside road in anavyssos bay

	ele, err = srtm.Get(gravelmap.Point{Lat: 37.7152154, Lng: 23.9347653})
	assert.Nil(t, err)
	assert.Equal(t, int32(0), ele) // Anavyssos bay point in the water
}
