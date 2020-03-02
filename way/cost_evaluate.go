package way

import (
	"github.com/thanosKontos/gravelmap"
)

const (
	// vehicleAcceptanceExclusively defines a way specifically made for the vehicle (e.g. cycleway for bicycles)
	vehicleAcceptanceExclusively = iota

	// vehicleAcceptanceNo defines a way that can be used by vehicle (a small city road for bicycles)
	vehicleAcceptanceYes

	// vehicleAcceptanceNo defines a way that can be used by vehicle but it not recommended (a larger road for bicycles)
	vehicleAcceptancePartially

	// vehicleAcceptanceNo defines a way that cannot be used by vehicle (e.g. footway for bicycles with no bike designation tags)
	vehicleAcceptanceMaybe

	// vehicleAcceptanceNo defines a way that cannot be used vehicle (e.g. path for SUVs)
	vehicleAcceptanceNo
)

const (
	// wayAcceptanceYes defines a way that is allowed to follow in a specific direction (e.g. a 2 way road)
	wayAcceptanceYes = iota

	// wayAcceptanceNo defines a way that is not allowed to follow in a specific direction (e.g. a direction off a one way road)
	wayAcceptanceNo
)

type wayAcceptance struct {
	normal int32
	reverse int32
}

type costEvaluate struct {
	distanceCalc gravelmap.DistanceCalculator
	elevationGetterCloser gravelmap.ElevationGetterCloser
}

func NewCostEvaluate(distanceCalc gravelmap.DistanceCalculator, elevationGetterCloser gravelmap.ElevationGetterCloser) *costEvaluate {
	return &costEvaluate{
		distanceCalc: distanceCalc,
		elevationGetterCloser: elevationGetterCloser,
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
	}

	wayAcceptance := getWayAcceptance(tags)
	vehicleAcceptance := getVehicleWayAcceptance(tags)

	wayAcceptanceWeightNormal := 1.0
	wayAcceptanceWeightReverse := 1.0
	if wayAcceptance.normal == wayAcceptanceNo {
		wayAcceptanceWeightNormal = 10000
	}
	if wayAcceptance.reverse == wayAcceptanceNo {
		wayAcceptanceWeightReverse = 10000
	}

	vehicleAcceptanceWeight := 1.0
	switch vehicleAcceptance {
	case vehicleAcceptanceExclusively:
		vehicleAcceptanceWeight = 0.7
	case vehicleAcceptancePartially:
		vehicleAcceptanceWeight = 2
	case vehicleAcceptanceMaybe:
		vehicleAcceptanceWeight = 1000
	case vehicleAcceptanceNo:
		vehicleAcceptanceWeight = 10000
	}

	offRoadWeight := 1.0
	wayType := gravelmap.WayTypePaved
	if isOffRoadWay(tags) {
		offRoadWeight = 0.6
		wayType = gravelmap.WayTypeUnaved
	}

	elevation, err := ce.elevationGetterCloser.Get(points, distance)
	elevationInfo := gravelmap.ElevationInfo{}
	if elevation != nil {
		elevationInfo = elevation.ElevationInfo
	}

	elevationWeight := elevationWeight{1, 1}
	if err == nil {
		elevationWeight = getElevationWeight(*elevation)
	}

	return gravelmap.WayEvaluation{
		Cost: int64(distance * vehicleAcceptanceWeight * wayAcceptanceWeightNormal * offRoadWeight * elevationWeight.weight),
		ReverseCost: int64(distance * vehicleAcceptanceWeight * wayAcceptanceWeightReverse * offRoadWeight * elevationWeight.reverse),
		Distance: int32(distance),
		WayType: wayType,
		ElevationInfo: elevationInfo,
	}
}

func getWayAcceptance(tags map[string]string) wayAcceptance {
	if val, ok := tags["oneway"]; ok {
		if val == "yes" {
			if val, ok := tags["cycleway"]; ok {
				if val == "opposite" || val == "opposite_lane" {
					return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
				}
			}

			if val, ok := tags["cycleway:left"]; ok {
				if val == "opposite_lane" {
					return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
				}
			}

			if val, ok := tags["cycleway:right"]; ok {
				if val == "opposite_lane" {
					return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
				}
			}

			if val, ok := tags["oneway:bicycle"]; ok {
				if val == "no" {
					return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
				}
			}
		}

		return wayAcceptance{wayAcceptanceNo, wayAcceptanceYes}
	}

	return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
}

func getVehicleWayAcceptance(tags map[string]string) int32 {
	if _, ok := tags["mtb:scale"]; ok {
		return vehicleAcceptanceExclusively
	}

	if val, ok := tags["bicycle"]; ok {
		if val == "yes" || val == "permissive" || val == "designated" {
			return vehicleAcceptanceExclusively
		}

		if val == "no" {
			return vehicleAcceptanceNo
		}
	}

	if val, ok := tags["highway"]; ok {
		if val == "footway" || val == "path" {
			return vehicleAcceptanceMaybe
		}

		if val == "cycleway" {
			return vehicleAcceptanceExclusively
		}

		if val == "service" {
			if val, ok := tags["bicycle"]; ok {
				if val == "yes" || val == "permissive" || val == "designated" {
					return vehicleAcceptanceExclusively
				}
			} else {
				return vehicleAcceptanceNo
			}
		}

		if val == "motorway" {
			return vehicleAcceptanceNo
		}

		if val == "primary" {
			return vehicleAcceptancePartially
		}

		return vehicleAcceptanceYes
	}

	return vehicleAcceptanceMaybe
}

func isOffRoadWay(tags map[string]string) bool {
	if val, ok := tags["surface"]; ok {
		if val == "unpaved" || val == "fine_gravel" || val == "gravel" || val == "compacted" || val == "pebblestone" || val == "earth" || val == "dirt" || val == "grass" || val == "ground" {
			return true
		}
	}

	if val, ok := tags["highway"]; ok {
		if val == "track" {
			return true
		}
	}

	return false
}

type elevationWeight struct {
	weight float64
	reverse float64
}

//0%: A flat road
//1-3%: Slightly uphill but not particularly challenging. A bit like riding into the wind.
//4-6%: A manageable gradient that can cause fatigue over long periods.
//7-9%: Starting to become uncomfortable for seasoned riders, and very challenging for new climbers.
//10%-15%: A painful gradient, especially if maintained for any length of time
//16%+: Very challenging for riders of all abilities. Maintaining this sort of incline for any length of time is very painful.
func getElevationWeight(elevation gravelmap.WayElevation) elevationWeight {
	switch {
	case elevation.Grade < -15:
		return elevationWeight{1, 15}
	case elevation.Grade < -10:
		return elevationWeight{1, 10}
	case elevation.Grade < -7:
		return elevationWeight{1, 7}
	case elevation.Grade < -4:
		return elevationWeight{0.8, 3}
	case elevation.Grade < -2:
		return elevationWeight{0.8, 1.2}
	case elevation.Grade < 0:
		return elevationWeight{0.8, 1}
	case elevation.Grade < 2:
		return elevationWeight{1, 0.8}
	case elevation.Grade < 4:
		return elevationWeight{1.2, 0.8}
	case elevation.Grade < 7:
		return elevationWeight{3, 0.8}
	case elevation.Grade < 10:
		return elevationWeight{7, 1}
	case elevation.Grade < 15:
		return elevationWeight{10, 1}
	default:
		return elevationWeight{15, 1}
	}
}