package path

import (
	"github.com/thanosKontos/gravelmap"
)

type simpleSimplifiedPath struct {
	distanceCalc gravelmap.DistanceCalculator
}

func NewSimpleSimplifiedPath(distanceCalc gravelmap.DistanceCalculator) *simpleSimplifiedPath {
	return &simpleSimplifiedPath{
		distanceCalc: distanceCalc,
	}
}

func (dp *simpleSimplifiedPath) Simplify(points []gravelmap.Point) []gravelmap.Point {
	if len(points) < 3 {
		return points
	}

	simplifiedPoints := []gravelmap.Point{points[0]}

	fromIdx := 0
	throughIdx := 1
	toIdx := 2

	for {
		if toIdx == len(points) {
			simplifiedPoints = append(simplifiedPoints, points[toIdx-1])

			break
		}

		d1 := dp.distanceCalc.Calculate(points[fromIdx], points[throughIdx])
		d2 := dp.distanceCalc.Calculate(points[throughIdx], points[toIdx])
		d3 := dp.distanceCalc.Calculate(points[fromIdx], points[toIdx])

		e := float64(d1+d2) / float64(d3)

		if e == 1 {
			throughIdx = toIdx
			toIdx++
		} else {
			simplifiedPoints = append(simplifiedPoints, points[throughIdx])

			fromIdx++
			throughIdx++
			toIdx++
		}
	}

	return simplifiedPoints
}
