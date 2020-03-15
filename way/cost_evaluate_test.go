package way

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/gravelmapfakes"
)

func TestEvaluateHappyPath(t *testing.T) {
	weigherFake := &gravelmapfakes.FakeWeighter{}
	elevationGetterCloserFake := &gravelmapfakes.FakeElevationGetterCloser{}
	distanceCalcFake := &gravelmapfakes.FakeDistanceCalculator{}

	normalEle := gravelmap.ElevationInfo{Grade: 5.5, From: 20, To: 50}
	reverseEle := gravelmap.ElevationInfo{Grade: -5.5, From: 50, To: 20}
	wayEle := gravelmap.WayElevation{
		ElevationEvaluation: gravelmap.ElevationEvaluation{
			Normal:  normalEle,
			Reverse: reverseEle,
		},
	}
	elevationGetterCloserFake.GetReturns(&wayEle, nil)

	distanceCalcFake.CalculateReturns(100)

	weigherFake.WeightOffRoadReturns(0.6)
	weigherFake.WeightWayAcceptanceReturns(gravelmap.Weight{0.6, 0.8})
	weigherFake.WeightVehicleAcceptanceReturns(0.8)
	weigherFake.WeightElevationReturns(gravelmap.Weight{0.5, 2})

	points := []gravelmap.Point{{38.0074488,23.8017056}, {38.0059424,23.7964307}, {38.0036874,23.7930924}}

	costEval := NewCostEvaluate(distanceCalcFake, elevationGetterCloserFake, weigherFake)
	tags := map[string]string{
		"surface": "asphalt",
		"highway": "track",
	}
	wayEval := costEval.Evaluate(points, tags)

	assert.Equal(t, normalEle, wayEval.ElevationEvaluation.Normal)
	assert.Equal(t, reverseEle, wayEval.ElevationEvaluation.Reverse)
	assert.Equal(t, int32(200), wayEval.Distance)
	assert.Equal(t, gravelmap.WayTypeUnaved, wayEval.WayType)
	assert.Equal(t, gravelmap.WayCost{28, 153}, wayEval.WayCost)
}

func TestEvaluateNoElevationData(t *testing.T) {
	weigherFake := &gravelmapfakes.FakeWeighter{}
	elevationGetterCloserFake := &gravelmapfakes.FakeElevationGetterCloser{}
	distanceCalcFake := &gravelmapfakes.FakeDistanceCalculator{}

	elevationGetterCloserFake.GetReturns(nil, errors.New("some error"))

	distanceCalcFake.CalculateReturns(100)

	weigherFake.WeightOffRoadReturns(0.6)
	weigherFake.WeightWayAcceptanceReturns(gravelmap.Weight{0.6, 0.8})
	weigherFake.WeightVehicleAcceptanceReturns(0.8)
	weigherFake.WeightElevationReturns(gravelmap.Weight{0.5, 2})

	points := []gravelmap.Point{{38.0074488,23.8017056}, {38.0059424,23.7964307}, {38.0036874,23.7930924}}

	costEval := NewCostEvaluate(distanceCalcFake, elevationGetterCloserFake, weigherFake)
	tags := map[string]string{
		"surface": "asphalt",
		"highway": "track",
	}
	wayEval := costEval.Evaluate(points, tags)

	assert.Equal(t, gravelmap.ElevationInfo{0.0, 0 , 0}, wayEval.ElevationEvaluation.Normal)
	assert.Equal(t, gravelmap.ElevationInfo{0.0, 0 , 0}, wayEval.ElevationEvaluation.Reverse)
	assert.Equal(t, int32(200), wayEval.Distance)
	assert.Equal(t, gravelmap.WayTypeUnaved, wayEval.WayType)
	assert.Equal(t, gravelmap.WayCost{28, 153}, wayEval.WayCost)
}
