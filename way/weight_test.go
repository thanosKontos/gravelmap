package way

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
)

func TestVehicleAcceptanceWeight(t *testing.T) {
	mtbConf := WeightConfig{}
	mtbConf.WeightVehicleAcceptance.Exclusively = 0.7
	mtbConf.WeightVehicleAcceptance.Yes = 1
	mtbConf.WeightVehicleAcceptance.Partially = 2
	mtbConf.WeightVehicleAcceptance.Maybe = 5
	mtbConf.WeightVehicleAcceptance.No = 10
	mtbConf.VehicleAcceptanceTags.Exclusively = TagValueConfig{
		Tags: []string{"mtb:scale"},
		Values: []map[string][]string{
			{"bicycle": []string{"yes", "permissive", "designated"}},
			{"highway": []string{"cycleway"}},
		},
	}
	mtbConf.VehicleAcceptanceTags.Maybe = TagValueConfig{
		Values: []map[string][]string{
			{"highway": []string{"path"}},
		},
	}
	mtbConf.VehicleAcceptanceTags.No = TagValueConfig{
		Values: []map[string][]string{
			{"bicycle": []string{"no"}},
		},
	}

	tests := map[string]struct {
		tags     map[string]string
		expected float64
	}{
		"residential_street_allowing_foot_and_bikes": {
			tags:     map[string]string{"foot": "yes", "highway": "cycleway"},
			expected: 0.7,
		},
		"path_with_mtb_note": {
			tags:     map[string]string{"highway": "path", "mtb:scale": "abcd1"},
			expected: 0.7,
		},
		"simple_path_with_no_other_info": {
			tags:     map[string]string{"highway": "path"},
			expected: 5,
		},
		"path_not_allowing_bikes": {
			tags:     map[string]string{"highway": "path", "bicycle": "no"},
			expected: 10,
		},
	}

	weighter := NewDefaultWeight(mtbConf)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, weighter.WeightVehicleAcceptance(test.tags))
		})
	}
}

func TestWayAcceptanceWeight(t *testing.T) {
	mtbConf := WeightConfig{}
	mtbConf.WayAcceptanceTags.Simple.NoDirection = TagValueConfig{
		Tags: []string{"military"},
	}
	mtbConf.WayAcceptanceTags.Simple.OppositeDirection = TagValueConfig{
		Values: []map[string][]string{
			{"oneway": []string{"yes"}},
		},
	}
	mtbConf.WayAcceptanceTags.Nested.BothDirection = append(mtbConf.WayAcceptanceTags.Nested.BothDirection, NestedTagValueConfig{
		Tag:          "oneway",
		Value:        "yes",
		NestedTag:    "cycleway",
		NestedValues: []string{"opposite", "opposite_lane"},
	})

	tests := map[string]struct {
		tags     map[string]string
		expected wayAcceptance
	}{
		"accessable_street": {
			tags:     map[string]string{"foot": "yes", "highway": "primary"},
			expected: wayAcceptance{wayAcceptanceYes, wayAcceptanceYes},
		},
		"unaccessable_street": {
			tags:     map[string]string{"military": "yes", "highway": "primary"},
			expected: wayAcceptance{wayAcceptanceNo, wayAcceptanceNo},
		},
		"one_way_street": {
			tags:     map[string]string{"highway": "primary", "oneway": "yes"},
			expected: wayAcceptance{wayAcceptanceYes, wayAcceptanceNo},
		},
		"one_way_street_with_bikepath": { // contradicting_simple_and_nested_wins_the_nested
			tags:     map[string]string{"highway": "primary", "oneway": "yes", "cycleway": "opposite"},
			expected: wayAcceptance{wayAcceptanceYes, wayAcceptanceYes},
		},
	}

	weighter := NewDefaultWeight(mtbConf)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, weighter.getMtbWayAcceptance(test.tags))
		})
	}
}

func TestElevationWeight(t *testing.T) {
	mtbConf := WeightConfig{}
	mtbConf.WeightElevation.Undefined = []float64{1.2, 1.2}
	mtbConf.WeightElevation.LessThan = map[float32][]float64{
		-12:  []float64{1, 15},
		-5:   []float64{1, 5},
		0:    []float64{1, 1},
		8:    []float64{6, 1},
		1000: []float64{18, 1},
	}

	tests := map[string]struct {
		ele      *gravelmap.WayElevation
		expected gravelmap.BidirectionalWeight
	}{
		"undefined_elevation_gives_default": {
			ele:      nil,
			expected: gravelmap.BidirectionalWeight{1.2, 1.2},
		},
		"min": {
			ele: &gravelmap.WayElevation{BidirectionalElevationInfo: gravelmap.BidirectionalElevationInfo{
				Normal:  gravelmap.ElevationInfo{Grade: -15},
				Reverse: gravelmap.ElevationInfo{Grade: 15},
			}},
			expected: gravelmap.BidirectionalWeight{1, 15},
		},
		"in_middle": {
			ele: &gravelmap.WayElevation{BidirectionalElevationInfo: gravelmap.BidirectionalElevationInfo{
				Normal:  gravelmap.ElevationInfo{Grade: 5.2},
				Reverse: gravelmap.ElevationInfo{Grade: -5.2},
			}},
			expected: gravelmap.BidirectionalWeight{6, 1},
		},
		"max": {
			ele: &gravelmap.WayElevation{BidirectionalElevationInfo: gravelmap.BidirectionalElevationInfo{
				Normal:  gravelmap.ElevationInfo{Grade: 15},
				Reverse: gravelmap.ElevationInfo{Grade: -15},
			}},
			expected: gravelmap.BidirectionalWeight{18, 1},
		},
	}

	weighter := NewDefaultWeight(mtbConf)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, weighter.WeightElevation(test.ele))
		})
	}
}
