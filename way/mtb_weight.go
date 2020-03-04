package way

import "github.com/thanosKontos/gravelmap"

type bicycleWeight struct {
}

func NewBicycleWeight() *bicycleWeight {
	return &bicycleWeight{}
}

func (b *bicycleWeight) WeightOffRoad(wayType int8) float64 {
	if wayType == gravelmap.WayTypeUnaved {
		return 0.6
	}

	return 1.0
}

func (b *bicycleWeight) WeightWayAcceptance(tags map[string]string) gravelmap.Weight {
	wayAcceptance := getWayAcceptance(tags)
	wayAcceptanceWeight := gravelmap.Weight{1.0, 1.0}
	if wayAcceptance.normal == wayAcceptanceNo {
		wayAcceptanceWeight.Normal = 10000
	}
	if wayAcceptance.reverse == wayAcceptanceNo {
		wayAcceptanceWeight.Reverse = 10000
	}

	return wayAcceptanceWeight
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

func (b *bicycleWeight) WeightVehicleAcceptance(tags map[string]string) float64 {
	switch getVehicleWayAcceptance(tags) {
	case vehicleAcceptanceExclusively:
		return 0.7
	case vehicleAcceptancePartially:
		return 2.0
	case vehicleAcceptanceMaybe:
	case vehicleAcceptanceNo:
		return 10000.0
	}

	return 1.0
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

//0%: A flat road
//1-3%: Slightly uphill but not particularly challenging. A bit like riding into the wind.
//4-6%: A manageable gradient that can cause fatigue over long periods.
//7-9%: Starting to become uncomfortable for seasoned riders, and very challenging for new climbers.
//10%-15%: A painful gradient, especially if maintained for any length of time
//16%+: Very challenging for riders of all abilities. Maintaining this sort of incline for any length of time is very painful.
func (b *bicycleWeight) WeightElevation(elevation *gravelmap.WayElevation) gravelmap.Weight {
	if elevation == nil {
		return gravelmap.Weight{1, 15}
	}

	switch {
	case elevation.ElevationEvaluation.Normal.Grade < -15:
		return gravelmap.Weight{1, 15}
	case elevation.ElevationEvaluation.Normal.Grade < -10:
		return gravelmap.Weight{1, 10}
	case elevation.ElevationEvaluation.Normal.Grade < -7:
		return gravelmap.Weight{1, 7}
	case elevation.ElevationEvaluation.Normal.Grade < -4:
		return gravelmap.Weight{0.8, 3}
	case elevation.ElevationEvaluation.Normal.Grade < -2:
		return gravelmap.Weight{0.8, 1.2}
	case elevation.ElevationEvaluation.Normal.Grade < 0:
		return gravelmap.Weight{0.8, 1}
	case elevation.ElevationEvaluation.Normal.Grade < 2:
		return gravelmap.Weight{1, 0.8}
	case elevation.ElevationEvaluation.Normal.Grade < 4:
		return gravelmap.Weight{1.2, 0.8}
	case elevation.ElevationEvaluation.Normal.Grade < 7:
		return gravelmap.Weight{3, 0.8}
	case elevation.ElevationEvaluation.Normal.Grade < 10:
		return gravelmap.Weight{7, 1}
	case elevation.ElevationEvaluation.Normal.Grade < 15:
		return gravelmap.Weight{10, 1}
	default:
		return gravelmap.Weight{15, 1}
	}
}
