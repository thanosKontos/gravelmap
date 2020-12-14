package way

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeightFromConfigAndTags(t *testing.T) {
	mtbConf := WeightConfig{
		WeightOffroad: 0.6,
	}
	mtbConf.WeightVehicleAcceptance.Exclusively = 0.7
	mtbConf.WeightVehicleAcceptance.Yes = 1
	mtbConf.WeightVehicleAcceptance.Partially = 2
	mtbConf.WeightVehicleAcceptance.Maybe = 10000000.0
	mtbConf.WeightVehicleAcceptance.No = 10000000.0
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

	tests := map[string]struct {
		tags                      map[string]string
		expectedVehicleAcceptance float64
	}{
		// "only_bike": {
		// 	tags:     map[string]string{"bicycle": "yes"},
		// 	expected: 0,
		// },
		// "simple_city_road": {
		// 	tags:     map[string]string{"highway": "residential", "surface": "asphalt"},
		// 	expected: 1,
		// },
		// "city_paved_road": {
		// 	tags:     map[string]string{"highway": "residential", "surface": "paving_stones"},
		// 	expected: 1,
		// },
		"residential_street_allowing_foot_and_bikes": {
			tags:                      map[string]string{"foot": "yes", "highway": "cycleway"},
			expectedVehicleAcceptance: 0.7,
		},
		"path_with_mtb_note": {
			tags:                      map[string]string{"highway": "path", "mtb:scale": "abcd1"},
			expectedVehicleAcceptance: 0.7,
		},
		"simple_path_with_no_other_info": {
			tags:                      map[string]string{"highway": "path"},
			expectedVehicleAcceptance: 10000000.0,
		},
		// "footway_not_mentioning_bike": {
		// 	tags:     map[string]string{"highway": "footway"},
		// 	expected: 3,
		// },
		// "no_tags": {
		// 	tags:     map[string]string{},
		// 	expected: 1,
		// },
	}

	weighter := NewBicycleWeight(mtbConf)

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expectedVehicleAcceptance, weighter.WeightVehicleAcceptance(test.tags))
		})
	}
}
