package srtm_ascii

import (
	"github.com/thanosKontos/gravelmap"
)

type ElevationGrader struct {
	elevationFinder gravelmap.ElevationFinder
}

// NewSRTM initialize and return an new SRTM object.
func NewElevationGrader(elevationFinder gravelmap.ElevationFinder) (*ElevationGrader, error) {
	return &ElevationGrader{
		elevationFinder: elevationFinder,
	}, nil
}

func (s *SRTM) Grade(points []gravelmap.Point) error {


	return nil
}
