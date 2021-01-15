package kml

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
)

func TestEmptyLegError(t *testing.T) {
	b := new(bytes.Buffer)
	routingLegs := []gravelmap.RoutingLeg{}

	kml := NewKml()
	err := kml.Write(b, routingLegs)

	assert.NotNil(t, err)
}

func TestKmlGeneration(t *testing.T) {
	b := new(bytes.Buffer)

	elevDecline := gravelmap.ElevationRange{140, 130}
	elevBigIncline := gravelmap.ElevationRange{130, 220}
	routingLegs := []gravelmap.RoutingLeg{
		gravelmap.RoutingLeg{
			Coordinates: []gravelmap.Point{gravelmap.Point{43.4, 22.4}, gravelmap.Point{45.4, 23.4}},
			Length:      12.6,
			WayType:     "unpaved",
			Elevation:   &elevDecline,
		},
		gravelmap.RoutingLeg{
			Coordinates: []gravelmap.Point{gravelmap.Point{44.4, 24.4}, gravelmap.Point{46.4, 25.4}},
			Length:      12.6,
			WayType:     "unpaved",
			Elevation:   &elevBigIncline,
		},
		gravelmap.RoutingLeg{
			Coordinates: []gravelmap.Point{gravelmap.Point{44.7, 24.7}, gravelmap.Point{46.7, 25.7}},
			Length:      15.6,
			WayType:     "unpaved",
			Elevation:   &elevBigIncline,
		},
	}

	kml := NewKml()
	err := kml.Write(b, routingLegs)

	assert.Nil(t, err)
	assert.Contains(t, b.String(), "<kml xmlns=\"http://www.opengis.net/kml/2.2\">")
	assert.Contains(t, b.String(), "<name>Extracted route from gravelmap</name>")
	assert.Contains(t, b.String(), " <description>Extracted route from gravelmap route from x to y.</description>")
	assert.Contains(t, b.String(), "<coordinates>22.4,43.4 23.4,45.4</coordinates>")
}
