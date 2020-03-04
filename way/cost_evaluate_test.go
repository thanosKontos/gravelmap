package way

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
