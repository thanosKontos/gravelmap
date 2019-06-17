package mock

import "github.com/thanosKontos/gravelmap"

type ElevationGrader struct {
	elevationFinder gravelmap.ElevationFinder
}

func NewElevationMock() *ElevationGrader {
	return &ElevationGrader{}
}

func (s *ElevationGrader) FindElevation(point gravelmap.Point) (float64, error) {
	return 10.0, nil
}
