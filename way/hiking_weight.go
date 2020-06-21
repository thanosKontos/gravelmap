package way

import (
	"github.com/thanosKontos/gravelmap"
	gmstring "github.com/thanosKontos/gravelmap/string"
)

type hikingWeight struct {
}

func NewHikingWeight() *hikingWeight {
	return &hikingWeight{}
}

func (b *hikingWeight) WeightOffRoad(wayType int8) float64 {
	if wayType == gravelmap.WayTypeUnaved {
		return 1.0
	}

	return 1.6
}

func (b *hikingWeight) WeightWayAcceptance(tags map[string]string) gravelmap.BidirectionalWeight {
	wayAcceptance := getFootWayAcceptance(tags)
	wayAcceptanceWeight := gravelmap.BidirectionalWeight{Normal: 1.0, Reverse: 1.0}
	if wayAcceptance.normal == wayAcceptanceNo {
		wayAcceptanceWeight.Normal = 10000000
	}
	if wayAcceptance.reverse == wayAcceptanceNo {
		wayAcceptanceWeight.Reverse = 10000000
	}

	return wayAcceptanceWeight
}

func getFootWayAcceptance(tags map[string]string) wayAcceptance {
	if val, ok := tags["highway"]; ok {
		if gmstring.String(val).Exists([]string{"motorway", "trunk", "primary"}) {
			return wayAcceptance{wayAcceptanceNo, wayAcceptanceNo}
		}
	}

	if _, ok := tags["military"]; ok {
		return wayAcceptance{wayAcceptanceNo, wayAcceptanceNo}
	}

	return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
}

func (b *hikingWeight) WeightVehicleAcceptance(tags map[string]string) float64 {
	switch getFootVehicleWayAcceptance(tags) {
	case vehicleAcceptanceExclusively:
		return 1.0
	case vehicleAcceptancePartially:
		return 2.2
	case vehicleAcceptanceMaybe:
		return 10000000.0
	case vehicleAcceptanceNo:
		return 10000000.0
	}

	return 1.0
}

func getFootVehicleWayAcceptance(tags map[string]string) int32 {
	if val, ok := tags["highway"]; ok {
		if gmstring.String(val).Exists([]string{"footway", "path", "pedestrian", "steps"}) {
			return vehicleAcceptanceExclusively
		}
	}

	return vehicleAcceptanceYes
}

//0%: A flat road
//1-3%: Slightly uphill but not particularly challenging. A bit like riding into the wind.
//4-6%: A manageable gradient that can cause fatigue over long periods.
//7-9%: Starting to become uncomfortable for seasoned riders, and very challenging for new climbers.
//10%-15%: A painful gradient, especially if maintained for any length of time
//16%+: Very challenging for riders of all abilities. Maintaining this sort of incline for any length of time is very painful.
func (b *hikingWeight) WeightElevation(elevation *gravelmap.WayElevation) gravelmap.BidirectionalWeight {
	if elevation == nil {
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 15}
	}

	switch {
	case elevation.BidirectionalElevationInfo.Normal.Grade < -15:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 6}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -10:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 4}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -7:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 2.5}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -4:
		return gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 1.5}
	case elevation.BidirectionalElevationInfo.Normal.Grade < -2:
		return gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 1.2}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 0:
		return gravelmap.BidirectionalWeight{Normal: 0.8, Reverse: 1}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 2:
		return gravelmap.BidirectionalWeight{Normal: 1, Reverse: 0.8}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 4:
		return gravelmap.BidirectionalWeight{Normal: 1.2, Reverse: 0.8}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 7:
		return gravelmap.BidirectionalWeight{Normal: 1.5, Reverse: 0.8}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 10:
		return gravelmap.BidirectionalWeight{Normal: 2.5, Reverse: 1}
	case elevation.BidirectionalElevationInfo.Normal.Grade < 15:
		return gravelmap.BidirectionalWeight{Normal: 4, Reverse: 1}
	default:
		return gravelmap.BidirectionalWeight{Normal: 6, Reverse: 1}
	}
}
