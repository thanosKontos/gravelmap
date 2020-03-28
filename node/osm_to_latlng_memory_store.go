package node

import (
	"errors"

	"github.com/thanosKontos/gravelmap"
)

type NodeLatLngMap map[int]gravelmap.Point

func NewOsm2LatLngMemoryStore() NodeLatLngMap {
	return make(NodeLatLngMap)
}

func (ms NodeLatLngMap) Write(osmID int, point gravelmap.Point) {
	ms[osmID] = point
}

func (ms NodeLatLngMap) Read(ndID int) (gravelmap.Point, error) {
	if val, ok := ms[ndID]; ok {
		return val, nil
	} else {
		return gravelmap.Point{}, errors.New("could not find latlng")
	}
}
