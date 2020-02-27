package way

import (
	"fmt"

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

type costEvaluate struct {
	distanceCalc gravelmap.DistanceCalculator
}

func NewCostEvaluate(distanceCalc gravelmap.DistanceCalculator) *costEvaluate {
	return &costEvaluate{
		distanceCalc: distanceCalc,
	}
}

func (ce *costEvaluate) Evaluate(way gravelmap.EvaluativeWay) gravelmap.WayCost {
	var distance float64 = 0.0
	prevPoint := gravelmap.Point{}
	for i, pt := range way.Points {
		if i == 0 {
			prevPoint = pt
			continue
		}

		distance += float64(ce.distanceCalc.Calculate(pt, prevPoint))
	}

	tagsWeight := getTagsWeight(way.Tags)

	//fmt.Println(tagsWeight)
	//fmt.Println("==================")

	return gravelmap.WayCost{
		Cost: int64(distance*tagsWeight.weight),
		ReverseCost: int64(distance*tagsWeight.oppositeWeight),
	}
}

type bidirectionalWeight struct {
	weight float64
	oppositeWeight float64
}

func getTagsWeight(tags map[string]string) bidirectionalWeight {
	fmt.Println(tags)

	// Remove tags with

	weight := bidirectionalWeight{1, 1}

	//if _, ok := tags["cycleway:right"]; ok {
	//	weight = 1
	//}
	//
	//if _, ok := tags["cycleway:left"]; ok {
	//	weight = 1
	//}

	if val, ok := tags["surface"]; ok {
		if val == "unpaved" || val == "gravel" || val == "fine_gravel" || val == "ground" {
			if val, ok := tags["bicycle"]; ok {
				if val == "yes" || val == "permissive" {
					weight = bidirectionalWeight{weight.weight*0.6, weight.oppositeWeight*0.6}
				}
			}
			weight = bidirectionalWeight{weight.weight*0.7, weight.oppositeWeight*0.7}
		}
	}

	if val, ok := tags["highway"]; ok {
		if val == "footway" || val == "path" {
			if val, ok := tags["bicycle"]; ok {
				if val == "yes" || val == "permissive" || val == "designated" {
					weight = bidirectionalWeight{weight.weight*0.7, weight.oppositeWeight*0.7}
				} else {
					weight = bidirectionalWeight{weight.weight*10000, weight.oppositeWeight*10000}
				}
			} else {
				weight = bidirectionalWeight{weight.weight*10000, weight.oppositeWeight*10000}
			}
		}

		if val == "cycleway" {
			weight = bidirectionalWeight{weight.weight*0.9, weight.oppositeWeight*0.9}
		}

		if val == "service" {
			if val, ok := tags["bicycle"]; ok {
				if val == "yes" || val == "permissive" || val == "designated" {
					weight = bidirectionalWeight{weight.weight*0.9, weight.oppositeWeight*0.9}
				}
			} else {
				return bidirectionalWeight{10000.0, 10000.0}
			}
		}
	}

	if _, ok := tags["cycleway"]; ok {
		weight = bidirectionalWeight{weight.weight*0.9, weight.oppositeWeight*0.9}
	}

	if val, ok := tags["bicycle"]; ok {
		if val == "yes" || val == "permissive" || val == "designated" {
			weight = bidirectionalWeight{weight.weight*0.9, weight.oppositeWeight*0.9}
		}
	}

	//map[cycleway:opposite highway:residential maxspeed:30 name:Emtinghauser Weg oneway:yes oneway:bicycle:no source:yahoo images source:maxspeed:DE:zone30 surface:asphalt]

	if val, ok := tags["oneway"]; ok {
		if val == "yes" {
			if val, ok := tags["cycleway"]; ok {
				if val == "opposite" {
					weight = bidirectionalWeight{1000.0, weight.weight}
				}
			} else {
				weight = bidirectionalWeight{weight.weight, 10000.0}
			}
		}
	}

	return weight
}

func getVehicleRoadAcceptance(tags map[string]string) int32 {
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

		return vehicleAcceptanceYes
	}

	return vehicleAcceptanceMaybe
}


