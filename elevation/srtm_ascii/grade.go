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

func (g *ElevationGrader) Grade(points []gravelmap.Point, distance float64) (*gravelmap.WayElevationOld, error) {
	if len(points) < 2 {
		return nil, errors.New("cannot grade way of one or less points")
	}

	prevPoint := points[0]
	prevElev := 0.0
	elevDiff := 0.0
	startElev := 0.0
	endElev := 0.0

	for i, point := range points {
		elev, err := g.elevationFinder.FindElevation(point)
		if err != nil {
			return nil, errors.New("could not grade because of missing point elevation")
		}

		if i == 0 {
			startElev = elev
		}
		endElev = elev

		if point != prevPoint {
			elevDiff += elev - prevElev
		}

		prevElev = elev
	}

	grade := 100 * elevDiff / distance
	if math.IsNaN(grade) {
		return nil, errors.New("division by zero")
	}

	return &gravelmap.WayElevationOld{
		Grade:  grade,
		Start:  startElev,
		End:    endElev,
		Length: distance,
	}, nil
}
