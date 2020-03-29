package way

import (
	"github.com/thanosKontos/gravelmap"
)

type costEvaluate struct {
	distanceCalc          gravelmap.DistanceCalculator
	elevationGetterCloser gravelmap.ElevationGetterCloser
	weighter              gravelmap.Weighter
}

func NewCostEvaluate(
	distanceCalc gravelmap.DistanceCalculator,
	elevationGetterCloser gravelmap.ElevationGetterCloser,
	weighter gravelmap.Weighter,
) *costEvaluate {
	return &costEvaluate{
		distanceCalc:          distanceCalc,
		elevationGetterCloser: elevationGetterCloser,
		weighter:              weighter,
	}
}

func (ce *costEvaluate) Evaluate(points []gravelmap.Point, tags map[string]string) gravelmap.WayEvaluation {
	var distance = 0.0
	prevPoint := gravelmap.Point{}

	for i, pt := range points {
		if i == 0 {
			prevPoint = pt
			continue
		}

		distance += float64(ce.distanceCalc.Calculate(pt, prevPoint))
		prevPoint = pt
	}

	wayType := gravelmap.WayTypePaved
	if isOffRoadWay(tags) {
		wayType = gravelmap.WayTypeUnaved
	}

	elevation, err := ce.elevationGetterCloser.Get(points, distance)
	elevationEval := gravelmap.BidirectionalElevationInfo{}
	if err == nil {
		elevationEval = elevation.BidirectionalElevationInfo
	}

	wayAcceptanceWeight := ce.weighter.WeightWayAcceptance(tags)
	vehicleAcceptanceWeight := ce.weighter.WeightVehicleAcceptance(tags)
	offRoadWeight := ce.weighter.WeightOffRoad(wayType)
	elevationWeight := ce.weighter.WeightElevation(elevation)

	return gravelmap.WayEvaluation{
		BidirectionalCost: gravelmap.BidirectionalCost{
			Normal:  int64(distance * vehicleAcceptanceWeight * wayAcceptanceWeight.Normal * offRoadWeight * elevationWeight.Normal),
			Reverse: int64(distance * vehicleAcceptanceWeight * wayAcceptanceWeight.Reverse * offRoadWeight * elevationWeight.Reverse),
		},
		Distance:                   int32(distance),
		WayType:                    wayType,
		BidirectionalElevationInfo: elevationEval,
	}
}

func isOffRoadWay(tags map[string]string) bool {
	if val, ok := tags["surface"]; ok {
		if val == "unpaved" || val == "fine_gravel" || val == "gravel" || val == "compacted" || val == "pebblestone" || val == "earth" || val == "dirt" || val == "grass" || val == "ground" {
			return true
		}
	}

	if val, ok := tags["highway"]; ok {
		if val == "track" {
			if val, ok := tags["surface"]; ok {
				if val == "paved" {
					return false
				}
			}

			return true
		}
	}

	return false
}
