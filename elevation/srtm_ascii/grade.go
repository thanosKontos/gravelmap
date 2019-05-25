package srtm_ascii

import (
	"errors"
	"github.com/thanosKontos/gravelmap"
)

type ElevationGrader struct {
	elevationFinder gravelmap.ElevationFinder
	distanceFinder  gravelmap.DistanceFinder
}

// NewElevationGrader initialize and return an new Elevation Grader object.
func NewElevationGrader(elevationFinder gravelmap.ElevationFinder, distanceFinder gravelmap.DistanceFinder) (*ElevationGrader, error) {
	return &ElevationGrader{
		elevationFinder: elevationFinder,
		distanceFinder:  distanceFinder,
	}, nil
}

func (g *ElevationGrader) Grade(points []gravelmap.Point) (float64, error) {
	if len(points) < 2 {
		return 0.0, errors.New("cannot grade way of one or less points")
	}

	prevPoint := points[0]
	prevElev := 0.0
	overallDistance := 0.0
	elevDiff := 0.0
	for _, point := range points {
		elev, err := g.elevationFinder.FindElevation(point)
		if err != nil {
			return 0.0, errors.New("could not grade because of missing point elevation")
		}

		if point != prevPoint {
			distance, _ := g.distanceFinder.Distance(prevPoint, point)
			overallDistance += distance
			elevDiff += elev - prevElev
		}

		prevElev = elev
	}

	return 100 * elevDiff / overallDistance, nil
}
