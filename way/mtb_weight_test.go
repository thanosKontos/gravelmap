package way

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
)

func TestGetWeightWayAcceptance(t *testing.T) {
	weight := NewBicycleWeight()

	assert.Equal(
		t,
		gravelmap.Weight{1, 1},
		weight.WeightWayAcceptance(map[string]string{"highway": "residential", "surface": "asphalt"}),
	)
	assert.Equal(
		t,
		gravelmap.Weight{1, 10000000},
		weight.WeightWayAcceptance(map[string]string{"oneway": "yes", "surface": "asphalt"}),
	)
	assert.Equal(
		t,
		gravelmap.Weight{1, 1},
		weight.WeightWayAcceptance(map[string]string{"oneway": "yes", "cycleway": "opposite"}),
	)
	assert.Equal(
		t,
		gravelmap.Weight{1, 10000000},
		weight.WeightWayAcceptance(map[string]string{"oneway": "yes", "cycleway": "opposite-lane"}),
	)
}

func TestGetWeightOffRoad(t *testing.T) {
	weight := NewBicycleWeight()

	assert.Equal(t, 0.6, weight.WeightOffRoad(gravelmap.WayTypeUnaved))
	assert.Equal(t, 1.0, weight.WeightOffRoad(gravelmap.WayTypePaved))
}

func TestGetTagsWeight(t *testing.T) {
	tests := map[string]struct {
		tags     map[string]string
		expected int32
	}{
		"only_bike": {
			tags:     map[string]string{"bicycle": "yes"},
			expected: 0,
		},
		"simple_city_road": {
			tags:     map[string]string{"highway": "residential", "surface": "asphalt"},
			expected: 1,
		},
		"city_paved_road": {
			tags:     map[string]string{"highway": "residential", "surface": "paving_stones"},
			expected: 1,
		},
		"cycleway_allowing_foot": {
			tags:     map[string]string{"foot": "yes", "highway": "cycleway"},
			expected: 0,
		},
		"footway_not_mentioning_bike": {
			tags:     map[string]string{"highway": "footway"},
			expected: 3,
		},
		"no_tags": {
			tags:     map[string]string{},
			expected: 1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			weight := getVehicleWayAcceptance(test.tags)
			assert.Equal(t, test.expected, weight)
		})
	}
}
