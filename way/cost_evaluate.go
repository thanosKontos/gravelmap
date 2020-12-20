package way

import (
	"github.com/thanosKontos/gravelmap"
	gmstring "github.com/thanosKontos/gravelmap/string"
)

type costEvaluate struct {
	distanceCalc             gravelmap.DistanceCalculator
	elevationWayGetterCloser gravelmap.ElevationWayGetterCloser
	weighter                 gravelmap.Weighter
}

func NewCostEvaluate(
	distanceCalc gravelmap.DistanceCalculator,
	elevationWayGetterCloser gravelmap.ElevationWayGetterCloser,
	weighter gravelmap.Weighter,
) *costEvaluate {
	return &costEvaluate{
		distanceCalc:             distanceCalc,
		elevationWayGetterCloser: elevationWayGetterCloser,
		weighter:                 weighter,
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
		wayType = gravelmap.WayTypeUnpaved
	}
	if isPathway(tags) {
		wayType = gravelmap.WayTypePath
	}

	elevation, err := ce.elevationWayGetterCloser.Get(points, distance)
	elevationInfo := gravelmap.ElevationInfo{}
	if err == nil {
		elevationInfo = elevation.ElevationInfo
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
		Distance:      int32(distance),
		WayType:       wayType,
		ElevationInfo: elevationInfo,
	}
}

func isOffRoadWay(tags map[string]string) bool {
	if val, ok := tags["surface"]; ok {
		if gmstring.String(val).Exists([]string{"unpaved", "fine_gravel", "gravel", "compacted", "pebblestone", "earth", "dirt", "grass", "ground"}) {
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

func isPathway(tags map[string]string) bool {
	if val, ok := tags["highway"]; ok {
		return gmstring.String(val).Exists([]string{"path", "pedestrian", "steps", "footway"})
	}

	return false
}
