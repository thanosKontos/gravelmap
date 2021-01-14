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

	er := gravelmap.ElevationRange{130, 140}
	routingLegs := []gravelmap.RoutingLeg{
		gravelmap.RoutingLeg{
			Coordinates: []gravelmap.Point{gravelmap.Point{43.4, 22.4}, gravelmap.Point{45.4, 23.4}},
			Length:      12.6,
			WayType:     "unpaved",
			Elevation:   &er,
		},
	}

	kml := NewKml()
	err := kml.Write(b, routingLegs)

	assert.Nil(t, err)
	assert.Contains(t, b.String(), "<kml xmlns=\"http://www.opengis.net/kml/2.2\">")
	assert.Contains(t, b.String(), "<name>Extracted route from gravelmap</name>")
	assert.Contains(t, b.String(), " <description>Extracted route from gravelmap route from x to y.</description>")
	assert.Contains(t, b.String(), "<coordinates>22.400000,43.400000,0")
	assert.Contains(t, b.String(), "23.400000,45.400000,0</coordinates>")
}
