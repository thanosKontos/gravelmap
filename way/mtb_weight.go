package way

import (
	"github.com/thanosKontos/gravelmap"
	gmstring "github.com/thanosKontos/gravelmap/string"
)

type bicycleWeight struct {
}

func NewBicycleWeight() *bicycleWeight {
	return &bicycleWeight{}
}

func (b *bicycleWeight) WeightOffRoad(wayType int8) float64 {
	if wayType == gravelmap.WayTypeUnpaved {
		return 1.0
	}

	return 1.6
}

func (b *bicycleWeight) WeightWayAcceptance(tags map[string]string) gravelmap.BidirectionalWeight {
	wayAcceptance := getMtbWayAcceptance(tags)
	wayAcceptanceWeight := gravelmap.BidirectionalWeight{Normal: 1.0, Reverse: 1.0}
	if wayAcceptance.normal == wayAcceptanceNo {
		wayAcceptanceWeight.Normal = 10000000
	}
	if wayAcceptance.reverse == wayAcceptanceNo {
		wayAcceptanceWeight.Reverse = 10000000
	}

	return wayAcceptanceWeight
}

func getMtbWayAcceptance(tags map[string]string) wayAcceptance {
	if val, ok := tags["oneway"]; ok {
		if val == "yes" {
			if val, ok := tags["cycleway"]; ok {
				if gmstring.String(val).Exists([]string{"opposite", "opposite_lane"}) {
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

			return wayAcceptance{wayAcceptanceYes, wayAcceptanceNo}
		}

		return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
	}

	if _, ok := tags["military"]; ok {
		return wayAcceptance{wayAcceptanceNo, wayAcceptanceNo}
	}

	return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
}

func (b *bicycleWeight) WeightVehicleAcceptance(tags map[string]string) float64 {
	switch getMtbVehicleWayAcceptance(tags) {
	case vehicleAcceptanceExclusively:
		return 0.7
	case vehicleAcceptancePartially:
		return 2.0
	case vehicleAcceptanceMaybe:
		return 10000000.0
	case vehicleAcceptanceNo:
		return 10000000.0
	}

	return 1.0
}

func getMtbVehicleWayAcceptance(tags map[string]string) int32 {
	if _, ok := tags["mtb:scale"]; ok {
		return vehicleAcceptanceExclusively
	}

	if val, ok := tags["bicycle"]; ok {
		if gmstring.String(val).Exists([]string{"yes", "permissive", "designated"}) {
			return vehicleAcceptanceExclusively
		}

		if val == "no" {
			return vehicleAcceptanceNo
		}
	}

	if val, ok := tags["highway"]; ok {
		if gmstring.String(val).Exists([]string{"footway", "path"}) {
			return vehicleAcceptanceMaybe
		}

		if val == "cycleway" {
			return vehicleAcceptanceExclusively
		}

		if val == "service" {
			if val, ok := tags["bicycle"]; ok {
				if gmstring.String(val).Exists([]string{"yes", "permissive", "designated"}) {
					return vehicleAcceptanceExclusively
				}
			} else {
				return vehicleAcceptanceNo
			}
		}

		if gmstring.String(val).Exists([]string{"motorway", "steps"}) {
			return vehicleAcceptanceNo
		}

		if val == "primary" {
			return vehicleAcceptancePartially
		}

		return vehicleAcceptanceYes
	}

	return vehicleAcceptanceYes
}

//0%: A flat road
//1-3%: Slightly uphill but not particularly challenging. A bit like riding into the wind.
//4-6%: A manageable gradient that can cause fatigue over long periods.
//7-9%: Starting to become uncomfortable for seasoned riders, and very challenging for new climbers.
//10%-15%: A painful gradient, especially if maintained for any length of time
//16%+: Very challenging for riders of all abilities. Maintaining this sort of incline for any length of time is very painful.
func (b *bicycleWeight) WeightElevation(elevation *gravelmap.WayElevation) gravelmap.BidirectionalWeight {
	if elevation == nil {
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 15}
	}

	switch {
	case elevation.BidirectionalElevationInfo.Normal.Grade < -15:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 15}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -10:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 10}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -7:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 7}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -4:
		return gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 3}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -2:
		return gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 1.2}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 0:
		return gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 1}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 2:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 0.8}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 4:
		return gravelmap.BidirectionalWeight{Normal: 1.2, Reverse: 0.8}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 7:
		return gravelmap.BidirectionalWeight{Normal: 3, Reverse: 0.8}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 10:
		return gravelmap.BidirectionalWeight{Normal: 7, Reverse: 1}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 15:
		return gravelmap.BidirectionalWeight{Normal: 10, Reverse: 1}
	default:
		return gravelmap.BidirectionalWeight{Normal: 15, Reverse: 1}
	}
}
