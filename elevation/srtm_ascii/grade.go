package srtm_ascii

import (
	"errors"
	"github.com/thanosKontos/gravelmap"
	"math"
)

type ElevationGrader struct {
	elevationFinder gravelmap.ElevationFinder
}

// NewElevationGrader initialize and return an new Elevation Grader object.
func NewElevationGrader(elevationFinder gravelmap.ElevationFinder) (*ElevationGrader, error) {
	return &ElevationGrader{
		elevationFinder: elevationFinder,
	}, nil
}

func (g *ElevationGrader) Grade(points []gravelmap.Point, distance float64) (float64, error) {
	if len(points) < 2 {
		return 0.0, errors.New("cannot grade way of one or less points")
	}

	prevPoint := points[0]
	prevElev := 0.0
	elevDiff := 0.0
	for _, point := range points {
		elev, err := g.elevationFinder.FindElevation(point)
		if err != nil {
			return 0.0, errors.New("could not grade because of missing point elevation")
		}

		if point != prevPoint {
			elevDiff += elev - prevElev
		}

		prevElev = elev
	}

	grade := 100 * elevDiff / distance
	if math.IsNaN(grade) {
		return 0.0, errors.New("division by zero")
	}

	return grade, nil
}
