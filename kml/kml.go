package kml

import (
	"errors"
	"fmt"
	"github.com/thanosKontos/gravelmap"
	"strings"
)

var kmlBase = `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Extracted route from gravelmap</name>
    <description>Extracted route from gravelmap route from x to y.</description>
	<Style id="black-pvd">
      <LineStyle>
        <color>ff000000</color>
        <width>4</width>
      </LineStyle>
    </Style>
	<Style id="black-upvd">
      <LineStyle>
        <color>ff000000</color>
        <width>2</width>
      </LineStyle>
    </Style>
    <Style id="red-pvd">
      <LineStyle>
        <color>ff5252ff</color>
        <width>4</width>
      </LineStyle>
    </Style>
    <Style id="red-upvd">
      <LineStyle>
        <color>ff5252ff</color>
        <width>2</width>
      </LineStyle>
    </Style>
    <Style id="green-pvd">
      <LineStyle>
        <color>7f236b31</color>
        <width>4</width>
      </LineStyle>
    </Style>
    <Style id="green-upvd">
      <LineStyle>
        <color>7f236b31</color>
        <width>2</width>
      </LineStyle>
    </Style>
    <Style id="pink-pvd">
      <LineStyle>
        <color>7fe863ff</color>
        <width>4</width>
      </LineStyle>
    </Style>
    <Style id="pink-upvd">
      <LineStyle>
        <color>7fe863ff</color>
        <width>2</width>
      </LineStyle>
    </Style>
    <Style id="blue-pvd">
      <LineStyle>
        <color>7f9e4a42</color>
        <width>4</width>
      </LineStyle>
    </Style>
    <Style id="blue-upvd">
      <LineStyle>
        <color>7f9e4a42</color>
        <width>2</width>
      </LineStyle>
    </Style>
    %s
  </Document>
</kml>
`

var placeMarkBase = `<Placemark>
<styleUrl>#%s</styleUrl>
<LineString>
<extrude>1</extrude>
<tessellate>1</tessellate>
<altitudeMode>absolute</altitudeMode>
<coordinates>%s</coordinates>
</LineString>
</Placemark>`

type kml struct{}

func NewKml() *kml {
	return &kml{}
}

func (k *kml) CreateFromRoute(routingLegs []gravelmap.RoutingLeg) (string, error) {
	if len(routingLegs) < 1 {
		return "", errors.New("not enough routing legs to create kml")
	}

	placemarks := ""
	for _, leg := range routingLegs {
		pointsSl := make([]string, 0)
		for _, point := range leg.Coordinates {
			pointsSl = append(pointsSl, fmt.Sprintf("%f,%f,0", point.Lng, point.Lat))
		}

		placemarkCoords := strings.Join(pointsSl, "\n")
		placemark := fmt.Sprintf(placeMarkBase, getToKmlLineColor(leg.Elevation, leg.Paved), placemarkCoords)

		placemarks += placemark
	}

	return fmt.Sprintf(kmlBase, placemarks), nil
}

func getToKmlLineColor(elevation *gravelmap.RoutingLegElevation, paved bool) string {
	surfaceStyle := "upvd"
	if paved {
		surfaceStyle = "pvd"
	}

	if elevation == nil {
		return "black-" + surfaceStyle
	}

	if elevation.Grade < 1 {
		return "green-" + surfaceStyle
	}

	if elevation.Grade < 1.3 {
		return "blue-" + surfaceStyle
	}

	if elevation.Grade < 3 {
		return "pink-" + surfaceStyle
	}

	return "red-" + surfaceStyle
}
