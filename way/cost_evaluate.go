package way

import "github.com/thanosKontos/gravelmap"

type costEvaluate struct {
	distanceCalc gravelmap.DistanceCalculator
}

func NewCostEvaluate(distanceCalc gravelmap.DistanceCalculator) *costEvaluate {
	return &costEvaluate{
		distanceCalc: distanceCalc,
	}
}

func (ce *costEvaluate) Evaluate(way gravelmap.EvaluativeWay) gravelmap.WayCost {
	var distance int64 = 0
	prevPoint := gravelmap.Point{}
	for i, pt := range way.Points {
		if i == 0 {
			prevPoint = pt
			continue
		}

		distance += ce.distanceCalc.Calculate(pt, prevPoint)
	}

	return gravelmap.WayCost{Cost: distance, ReverseCost: distance}
}
