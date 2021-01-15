package kml

import (
	"errors"
	"fmt"
	"image/color"
	"io"

	"github.com/thanosKontos/gravelmap"
	kml2 "github.com/twpayne/go-kml"
)

type kml struct{}

// NewKml instantiates a new kml object
func NewKml() *kml {
	return &kml{}
}

func (k *kml) Write(w io.Writer, routingLegs []gravelmap.RoutingLeg) error {
	if len(routingLegs) < 1 {
		return errors.New("not enough routing legs to create kml")
	}

	doc := kml2.Document(
		kml2.Name("Extracted route from gravelmap"),
		kml2.Description("Extracted route from gravelmap route from x to y."),
		kml2.SharedStyle("black-pvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 0, G: 0, B: 0, A: 255}),
			kml2.Width(4),
		)),
		kml2.SharedStyle("black-upvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 0, G: 0, B: 0, A: 255}),
			kml2.Width(2),
		)),
		kml2.SharedStyle("red-pvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 255, G: 82, B: 82, A: 255}),
			kml2.Width(4),
		)),
		kml2.SharedStyle("red-upvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 255, G: 82, B: 82, A: 255}),
			kml2.Width(2),
		)),
		kml2.SharedStyle("green-pvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 49, G: 107, B: 35, A: 127}),
			kml2.Width(4),
		)),
		kml2.SharedStyle("green-upvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 49, G: 107, B: 35, A: 127}),
			kml2.Width(2),
		)),
		kml2.SharedStyle("pink-pvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 255, G: 99, B: 232, A: 127}),
			kml2.Width(4),
		)),
		kml2.SharedStyle("pink-upvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 255, G: 99, B: 232, A: 127}),
			kml2.Width(2),
		)),
		kml2.SharedStyle("blue-pvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 0, G: 128, B: 128, A: 127}),
			kml2.Width(4),
		)),
		kml2.SharedStyle("blue-upvd", kml2.LineStyle(
			kml2.Color(color.RGBA{R: 0, G: 128, B: 128, A: 127}),
			kml2.Width(2),
		)),
	)

	for _, way := range routingLegs {

		var routeCoordinates []kml2.Coordinate = []kml2.Coordinate{}
		for _, point := range way.Coordinates {
			routeCoordinates = append(routeCoordinates, kml2.Coordinate{Lon: point.Lng, Lat: point.Lat})
		}

		pm := kml2.Placemark(
			kml2.StyleURL(fmt.Sprintf("#%s", getWayLineStyle(way))),
			kml2.LineString(
				kml2.Coordinates(routeCoordinates...),
				kml2.Tessellate(true),
			),
		)

		doc.Add(pm)
	}

	result := kml2.KML(doc)

	return result.WriteIndent(w, "", "  ")
}

func getWayLineStyle(way gravelmap.RoutingLeg) string {
	surfaceStyle := "pvd"
	if way.WayType == "unpaved" || way.WayType == "path" {
		surfaceStyle = "upvd"
	}

	if way.Elevation == nil {
		return "black-" + surfaceStyle
	}

	grade := (way.Elevation.End - way.Elevation.Start) * 100 / way.Length
	if grade < 1 {
		return "green-" + surfaceStyle
	}

	if grade < 1.3 {
		return "blue-" + surfaceStyle
	}

	if grade < 3 {
		return "pink-" + surfaceStyle
	}

	return "red-" + surfaceStyle
}
