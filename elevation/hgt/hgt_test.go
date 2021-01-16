package hgt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/log"
)

type mockElevationStorage struct {
}

type mockElevation struct {
}

func (mes mockElevationStorage) Get(dms string) (gravelmap.ElevationPointGetter, error) {
	return mockElevation{}, nil
}

func (mes mockElevationStorage) Close() {
}

func (me mockElevation) Get(pt gravelmap.Point) (int32, error) {
	return int32(pt.Lat), nil
}

func TestGetWayElevation(t *testing.T) {
	hgt := hgt{
		elevationFileStorage: mockElevationStorage{},
		logger:               log.NewNullLog(),
	}

	pts := []gravelmap.Point{gravelmap.Point{43.4, 22.4}, gravelmap.Point{45.4, 23.4}, gravelmap.Point{60.4, 25.4}}
	wayEle, err := hgt.Get(pts, 150.5)

	assert.Nil(t, err)
	assert.Equal(t, []int32{43, 45, 60}, wayEle.Elevations)
	assert.Equal(t, float32(11.295681), wayEle.ElevationInfo.Grade)
	assert.Equal(t, int16(43), wayEle.ElevationInfo.From)
	assert.Equal(t, int16(60), wayEle.ElevationInfo.To)
}

func TestSmallWayElevationGivesError(t *testing.T) {
	hgt := hgt{
		elevationFileStorage: mockElevationStorage{},
		logger:               log.NewNullLog(),
	}

	pts := []gravelmap.Point{gravelmap.Point{43.4, 22.4}, gravelmap.Point{45.4, 23.4}}
	_, err := hgt.Get(pts, 9.5)

	assert.NotNil(t, err)
}
