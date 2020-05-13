package way

import (
	"strconv"

	"github.com/thanosKontos/gravelmap"
)

type bicycleWeight struct {
}

func NewBicycleWeight() *bicycleWeight {
	return &bicycleWeight{}
}

func (b *bicycleWeight) WeightOffRoad(wayType int8) float64 {
	if wayType == gravelmap.WayTypeUnaved {
		return 0.6
	}

	return 5.0
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
	if val, ok := tags["mtb:scale:imba"]; ok {
		imbaScale, err := strconv.Atoi(val)
		if err == nil {
			if imbaScale >= 4 {
				return vehicleAcceptanceNo
			} else if imbaScale >= 2 {
				return vehicleAcceptanceYes
			} else {
				return vehicleAcceptanceExclusively
			}
		}
	}

	if val, ok := tags["mtb:scale"]; ok {
		mbScale, err := strconv.Atoi(val)
		if err == nil {
			if mbScale >= 3 {
				return vehicleAcceptanceNo
			} else if mbScale >= 2 {
				return vehicleAcceptanceYes
			} else {
				return vehicleAcceptanceExclusively
			}
		}
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

		if val == "motorway" || val == "steps" {
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
func (b *bicycleWeight) WeightElevation(tags map[string]string, elevation *gravelmap.WayElevation) gravelmap.BidirectionalWeight {
	elevWeight := gravelmap.BidirectionalWeight{Normal: 1, Reverse: 15}
	if elevation != nil {
		switch {
		case elevation.BidirectionalElevationInfo.Normal.Grade < -15:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 1, Reverse: 15}
		case elevation.BidirectionalElevationInfo.Normal.Grade < -10:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 1, Reverse: 10}
		case elevation.BidirectionalElevationInfo.Normal.Grade < -7:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 1, Reverse: 7}
		case elevation.BidirectionalElevationInfo.Normal.Grade < -4:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 3}
		case elevation.BidirectionalElevationInfo.Normal.Grade < -2:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 1.2}
		case elevation.BidirectionalElevationInfo.Normal.Grade < 0:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 1}
		case elevation.BidirectionalElevationInfo.Normal.Grade < 2:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 1, Reverse: 0.8}
		case elevation.BidirectionalElevationInfo.Normal.Grade < 4:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 1.2, Reverse: 0.8}
		case elevation.BidirectionalElevationInfo.Normal.Grade < 7:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 3, Reverse: 0.8}
		case elevation.BidirectionalElevationInfo.Normal.Grade < 10:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 7, Reverse: 1}
		case elevation.BidirectionalElevationInfo.Normal.Grade < 15:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 10, Reverse: 1}
		default:
			elevWeight = gravelmap.BidirectionalWeight{Normal: 15, Reverse: 1}
		}
	}

	if val, ok := tags["mtb:scale:uphill"]; ok {
		uphill, err := strconv.Atoi(val)
		if err == nil {
			if uphill >= 4 {
				elevWeight = gravelmap.BidirectionalWeight{Normal: 15, Reverse: 1}
			} else if uphill >= 2 {
				elevWeight = gravelmap.BidirectionalWeight{Normal: 10, Reverse: 1}
			} else {
				elevWeight = gravelmap.BidirectionalWeight{Normal: 3, Reverse: 1}
			}
		}
	}

	return elevWeight
}
