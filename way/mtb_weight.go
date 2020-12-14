package way

import (
	"io/ioutil"
	"log"
	"sort"

	"github.com/thanosKontos/gravelmap"
	gmstring "github.com/thanosKontos/gravelmap/string"
	"gopkg.in/yaml.v2"
)

// type TagValueConfig struct {
// 	Tags   []string
// 	Values []map[string][]string
// }

type WeightConfig struct {
	WeightOffroad           float64 `yaml:"weight_offroad"`
	WeightVehicleAcceptance struct {
		Exclusively float64
		Yes         float64
		Partially   float64
		Maybe       float64
		No          float64
	} `yaml:"weight_vehicle_acceptance"`

	VehicleAcceptanceTags struct {
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
	} `yaml:"vehicle_acceptance_tags"`

	WayAcceptanceTags struct {
		Simple struct {
			NoDirection struct {
				Tags   []string
				Values []map[string][]string
			} `yaml:"no_direction"`
			OppositeDirection struct {
				Tags   []string
				Values []map[string][]string
			} `yaml:"opposite_direction"`
		}
		Nested struct {
			BothDirection []struct {
				Tag          string
				Value        string
				NestedTag    string   `yaml:"nested_tag"`
				NestedValues []string `yaml:"nested_values"`
			} `yaml:"both_direction"`
		}
	} `yaml:"way_acceptance_tags"`

	WeightElevation struct {
		Undefined []float64
		LessThan  map[float32][]float64 `yaml:"less_than"`
	} `yaml:"weight_elevation"`
}

type bicycleWeight struct {
	conf WeightConfig
}

func NewBicycleWeight() *bicycleWeight {
	conf := WeightConfig{}

	yamlFile, kkkerr := ioutil.ReadFile("profiles/mtb.yaml")
	if kkkerr != nil {
		log.Fatalf("error: %v", kkkerr)
	}
	errzzz := yaml.Unmarshal(yamlFile, &conf)
	if errzzz != nil {
		log.Fatalf("error: %v", errzzz)
	}

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
	wayAcceptance := b.getMtbWayAcceptance(tags)
	wayAcceptanceWeight := gravelmap.BidirectionalWeight{Normal: 1.0, Reverse: 1.0}
	if wayAcceptance.normal == wayAcceptanceNo {
		wayAcceptanceWeight.Normal = 10000000
	}
	if wayAcceptance.reverse == wayAcceptanceNo {
		wayAcceptanceWeight.Reverse = 10000000
	}

	return wayAcceptanceWeight
}

func (b *bicycleWeight) getMtbWayAcceptance(tags map[string]string) wayAcceptance {
	// First evaluate the nested configs if any
	for _, nestedConf := range b.conf.WayAcceptanceTags.Nested.BothDirection {
		if val, ok := tags[nestedConf.Tag]; ok {
			if val == nestedConf.Value {
				if val, ok := tags[nestedConf.NestedTag]; ok {
					if gmstring.String(val).Exists(nestedConf.NestedValues) {
						return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
					}
				}
			}
		}
	}

	// And then evaluate the simple cases
	for _, tag := range b.conf.WayAcceptanceTags.Simple.NoDirection.Tags {
		if _, ok := tags[tag]; ok {
			return wayAcceptance{wayAcceptanceNo, wayAcceptanceNo}
		}
	}

	for _, v := range b.conf.WayAcceptanceTags.Simple.NoDirection.Values {
		for tag, vals := range v {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return wayAcceptance{wayAcceptanceNo, wayAcceptanceNo}
				}
			}
		}
	}

	for _, tag := range b.conf.WayAcceptanceTags.Simple.OppositeDirection.Tags {
		if _, ok := tags[tag]; ok {
			return wayAcceptance{wayAcceptanceYes, wayAcceptanceNo}
		}
	}

	for _, v := range b.conf.WayAcceptanceTags.Simple.OppositeDirection.Values {
		for tag, vals := range v {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return wayAcceptance{wayAcceptanceYes, wayAcceptanceNo}
				}
			}
		}
	}

	// Allowing the vehicle to travel both directions is the default
	return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}

	// if _, ok := tags["military"]; ok {
	// 	return wayAcceptance{wayAcceptanceNo, wayAcceptanceNo}
	// }

	// if val, ok := tags["oneway"]; ok {
	// 	if val == "yes" {
	// 		if val, ok := tags["cycleway"]; ok {
	// 			if gmstring.String(val).Exists([]string{"opposite", "opposite_lane"}) {
	// 				return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
	// 			}
	// 		}

	// 		if val, ok := tags["cycleway:left"]; ok {
	// 			if val == "opposite_lane" {
	// 				return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
	// 			}
	// 		}

	// 		if val, ok := tags["cycleway:right"]; ok {
	// 			if val == "opposite_lane" {
	// 				return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
	// 			}
	// 		}

	// 		if val, ok := tags["oneway:bicycle"]; ok {
	// 			if val == "no" {
	// 				return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
	// 			}
	// 		}

	// 		return wayAcceptance{wayAcceptanceYes, wayAcceptanceNo}
	// 	}

	// 	return wayAcceptance{wayAcceptanceYes, wayAcceptanceYes}
	// }
}

func (b *bicycleWeight) WeightVehicleAcceptance(tags map[string]string) float64 {
	switch b.getMtbVehicleWayAcceptance(tags) {
	case vehicleAcceptanceExclusively:
		return b.conf.WeightVehicleAcceptance.Exclusively
	case vehicleAcceptancePartially:
		return b.conf.WeightVehicleAcceptance.Partially
	case vehicleAcceptanceMaybe:
		return b.conf.WeightVehicleAcceptance.Maybe
	case vehicleAcceptanceNo:
		return b.conf.WeightVehicleAcceptance.No
	}

	return b.conf.WeightVehicleAcceptance.Yes
}

func (b *bicycleWeight) getMtbVehicleWayAcceptance(tags map[string]string) int32 {
	for _, excl_tag := range b.conf.VehicleAcceptanceTags.Exclusively.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptanceExclusively
		}
	}

	for _, excl_vals := range b.conf.VehicleAcceptanceTags.Exclusively.Values {
		for tag, vals := range excl_vals {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return vehicleAcceptanceExclusively
				}
			}
		}
	}

	for _, excl_tag := range b.conf.VehicleAcceptanceTags.No.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptanceNo
		}
	}

	for _, excl_vals := range b.conf.VehicleAcceptanceTags.No.Values {
		for tag, vals := range excl_vals {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return vehicleAcceptanceNo
				}
			}
		}
	}

	for _, excl_tag := range b.conf.VehicleAcceptanceTags.Maybe.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptanceMaybe
		}
	}

	for _, excl_vals := range b.conf.VehicleAcceptanceTags.Maybe.Values {
		for tag, vals := range excl_vals {
			if val, ok := tags[tag]; ok {
				if gmstring.String(val).Exists(vals) {
					return vehicleAcceptanceMaybe
				}
			}
		}
	}

	for _, excl_tag := range b.conf.VehicleAcceptanceTags.Partially.Tags {
		if _, ok := tags[excl_tag]; ok {
			return vehicleAcceptancePartially
		}
	}

	for _, excl_vals := range b.conf.VehicleAcceptanceTags.Partially.Values {
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
			Normal:  b.conf.WeightElevation.Undefined[0],
			Reverse: b.conf.WeightElevation.Undefined[1],
		}
	}

	// Looping through map may not keep the order, so we are sorting the keys and looping using the sorted array
	confGrades := make([]float32, 0)
	for confGrade, _ := range b.conf.WeightElevation.LessThan {
		confGrades = append(confGrades, confGrade)
	}
	sort.Slice(confGrades, func(i, j int) bool { return confGrades[i] < confGrades[j] })
	for _, confGrade := range confGrades {
		if elevation.BidirectionalElevationInfo.Normal.Grade < confGrade {
			return gravelmap.BidirectionalWeight{
				Normal:  b.conf.WeightElevation.LessThan[confGrade][0],
				Reverse: b.conf.WeightElevation.LessThan[confGrade][1],
			}
		}
	}

	return gravelmap.BidirectionalWeight{
		Normal:  b.conf.WeightElevation.Undefined[0],
		Reverse: b.conf.WeightElevation.Undefined[1],
	}
}
