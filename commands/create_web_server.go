package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/cli"
	"github.com/thanosKontos/gravelmap/kml"
	"github.com/thanosKontos/gravelmap/route"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var kmlBase = `<?xml version="1.0" encoding="UTF-8"?>
<kml xmlns="http://www.opengis.net/kml/2.2">
  <Document>
    <name>Extracted route from gravelmap</name>
    <description>Extracted route from gravelmap route from x to y.</description>
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

// createRoutingDataCommand defines the create route command.
func createServerCommand() *cobra.Command {
	createServerCmd := &cobra.Command{
		Use:   "create-server",
		Short: "create a simple server to host a test route website",
		Long:  "create a simple server to host a test route website",
	}

	createServerCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createServerCmdRun()
	}

	return createServerCmd
}

// createRoutingDataCmdRun defines the command run actions.
func createServerCmdRun() error {
	http.HandleFunc("/route", routeHandler)
	http.HandleFunc("/create-kml", createKmlHandler)

	http.ListenAndServe(":8000", nil)

	return nil
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	pointFrom, err := getPointFromParams("from", r)
	pointTo, err2 := getPointFromParams("to", r)
	if err != nil || err2 != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}

	pgRouter, err := route.NewPgRouting(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"), cli.NewCLI())
	if err != nil {
		log.Fatal(err)
	}

	features, err := pgRouter.Route(
		*pointFrom,
		*pointTo,
	)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		return
	}

	json, _ := json.Marshal(features)
	fmt.Fprintf(w, string(json))
}

func createKmlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pointFrom, err := getPointFromParams("from", r)
	pointTo, err2 := getPointFromParams("to", r)
	if err != nil || err2 != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}

	pgRouter, err := route.NewPgRouting(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"), cli.NewCLI())
	if err != nil {
		log.Fatal(err)
	}

	features, err := pgRouter.Route(
		*pointFrom,
		*pointTo,
	)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		return
	}

	kml := kml.NewKml()
	kmlString, err := kml.CreateFromRoute(features)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		return
	}

	w.Header().Set("Content-Type", "application/vnd.google-earth.kml+xml")
	w.Header().Set("Content-Disposition", "attachment; filename=\"test.kml\"")
	fmt.Fprintf(w, kmlString)
}

func getPointFromParams(param string, r *http.Request) (*gravelmap.Point, error) {
	fromKeys, ok := r.URL.Query()[param]
	if !ok || len(fromKeys[0]) < 1 {
		return nil, errors.New("non existing param")
	}
	latLng := strings.Split(fromKeys[0], ",")

	lat, err := strconv.ParseFloat(latLng[1], 64)
	if err != nil {
		return nil, err
	}

	lng, err := strconv.ParseFloat(latLng[0], 64)
	if err != nil {
		return nil, err
	}

	return &gravelmap.Point{Lat: lat, Lng: lng}, nil
}
