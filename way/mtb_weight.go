package way

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/thanosKontos/gravelmap"
	gmstring "github.com/thanosKontos/gravelmap/string"
	"gopkg.in/yaml.v2"
)

type WeightConfig struct {
	WeightOffroad             float64 `yaml:"weight_offroad"`
	Weight_vehicle_acceptance struct {
		Exclusively float64
		Yes         float64
		Partially   float64
		Maybe       float64
		No          float64
	}

	Vehicle_acceptance_tags struct {
		Exclusively struct {
			Tags   []string
			Values []map[string][]string
		}
		Yes struct {
			Tags   []string
			Values []map[string][]string
		}
		Partially struct {
			Tags   []string
			Values []map[string][]string
		}
		Maybe struct {
			Tags   []string
			Values []map[string][]string
		}
		No struct {
			Tags   []string
			Values []map[string][]string
		}
	}

	Way_acceptance_tags struct {
		Tags interface{}
	}

	Weight_elevation struct {
		Undefined []float64
		Less_than map[float32][]float64
	}
}

type bicycleWeight struct {
	conf WeightConfig
}

func NewBicycleWeight() *bicycleWeight {
	conf := WeightConfig{}

	yamlFile, kkkerr := ioutil.ReadFile("profiles/mtb.yaml")
	if kkkerr != nil {
		fmt.Println(kkkerr)
		os.Exit(0)
	}
	errzzz := yaml.Unmarshal(yamlFile, &conf)
	if errzzz != nil {
		log.Fatalf("error: %v", errzzz)
	}

	//fmt.Println(conf.Way_acceptance_tags.Tags)
	fmt.Println(conf.Weight_elevation)

	return &bicycleWeight{
		conf,
	}
}

func (b *bicycleWeight) WeightOffRoad(wayType int8) float64 {
	if wayType == gravelmap.WayTypeUnpaved {
		return b.conf.WeightOffroad
	}

	return 1.0
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
	if _, ok := tags["military"]; ok {
		return wayAcceptance{wayAcceptanceNo, wayAcceptanceNo}
	}

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

	return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
}

func (b *bicycleWeight) WeightVehicleAcceptance(tags map[string]string) float64 {
	switch b.getMtbVehicleWayAcceptance(tags) {
	case vehicleAcceptanceExclusively:
		return b.conf.Weight_vehicle_acceptance.Exclusively
	case vehicleAcceptancePartially:
		return b.conf.Weight_vehicle_acceptance.Partially
	case vehicleAcceptanceMaybe:
		return b.conf.Weight_vehicle_acceptance.Maybe
	case vehicleAcceptanceNo:
		return b.conf.Weight_vehicle_acceptance.No
	}

	return b.conf.Weight_vehicle_acceptance.Yes
}

func (b *bicycleWeight) getMtbVehicleWayAcceptance(tags map[string]string) int32 {
	for _, excl_tag := range b.conf.Vehicle_acceptance_tags.Exclusively.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptanceExclusively
		}
	}

	for _, excl_vals := range b.conf.Vehicle_acceptance_tags.Exclusively.Values {
		for tag, vals := range excl_vals {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return vehicleAcceptanceExclusively
				}
			}
		}
	}

	for _, excl_tag := range b.conf.Vehicle_acceptance_tags.No.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptanceNo
		}
	}

	for _, excl_vals := range b.conf.Vehicle_acceptance_tags.No.Values {
		for tag, vals := range excl_vals {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return vehicleAcceptanceNo
				}
			}
		}
	}

	for _, excl_tag := range b.conf.Vehicle_acceptance_tags.Maybe.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptanceMaybe
		}
	}

	for _, excl_vals := range b.conf.Vehicle_acceptance_tags.Maybe.Values {
		for tag, vals := range excl_vals {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return vehicleAcceptanceMaybe
				}
			}
		}
	}

	for _, excl_tag := range b.conf.Vehicle_acceptance_tags.Partially.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptancePartially
		}
	}

	for _, excl_vals := range b.conf.Vehicle_acceptance_tags.Partially.Values {
		for tag, vals := range excl_vals {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return vehicleAcceptancePartially
				}
			}
		}
	}

	return vehicleAcceptanceYes
}

func (b *bicycleWeight) WeightElevation(elevation *gravelmap.WayElevation) gravelmap.BidirectionalWeight {
	if elevation == nil {
		return gravelmap.BidirectionalWeight{
			Normal:  b.conf.Weight_elevation.Undefined[0],
			Reverse: b.conf.Weight_elevation.Undefined[1],
		}
	}

	// Looping through map may not keep the order, so we are sorting the keys and looping using the sorted array
	confGrades := make([]float32, 0)
	for confGrade, _ := range b.conf.Weight_elevation.Less_than {
		confGrades = append(confGrades, confGrade)
	}
	sort.Slice(confGrades, func(i, j int) bool { return confGrades[i] < confGrades[j] })
	for _, confGrade := range confGrades {
		if elevation.BidirectionalElevationInfo.Normal.Grade < confGrade {
			return gravelmap.BidirectionalWeight{
				Normal:  b.conf.Weight_elevation.Less_than[confGrade][0],
				Reverse: b.conf.Weight_elevation.Less_than[confGrade][1],
			}
		}
	}

	return gravelmap.BidirectionalWeight{
		Normal:  b.conf.Weight_elevation.Undefined[0],
		Reverse: b.conf.Weight_elevation.Undefined[1],
	}
}
