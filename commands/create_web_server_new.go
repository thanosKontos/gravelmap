package commands

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/dijkstra"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/edge"
	"github.com/thanosKontos/gravelmap/way"
	"googlemaps.github.io/maps"

	"net/http"
)

// createWebServerNewCommand defines the create server command.
func createWebServerNewCommand() *cobra.Command {
	createWebServerNewCmd := &cobra.Command{
		Use:   "create-web-server-new",
		Short: "create a simple server to host a test route website",
		Long:  "create a simple server to host a test route website",
	}

	createWebServerNewCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createWebServerNewCmdRun()
	}

	return createWebServerNewCmd
}

// createRoutingDataCmdRun defines the command run actions.
func createWebServerNewCmdRun() error {
	http.HandleFunc("/route", routeNewHandler)

	http.ListenAndServe(":8000", nil)

	return nil
}

func routeNewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	distanceCalc := distance.NewHaversine()
	bboxFr := edge.NewBBoxFileRead("_files", distanceCalc)

	pointFrom, err := getPointFromParamsNew("from", r)
	pointTo, err2 := getPointFromParamsNew("to", r)
	if err != nil || err2 != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}

	edgeFrom, err := bboxFr.FindClosest(*pointFrom)
	if err != nil {
		w.WriteHeader(400)
		json, _ := json.Marshal(err)
		fmt.Fprintf(w, string(json))

		return
	}

	edgeTo, err := bboxFr.FindClosest(*pointTo)
	if err != nil {
		w.WriteHeader(400)
		json, _ := json.Marshal(err)
		fmt.Fprintf(w, string(json))

		return
	}


	graph := dijkstra.NewGraph()
	dataFile, err := os.Open("_files/graph.gob")
	if err != nil {
		fmt.Fprintf(w, `{"message": "Cannot open graph file"}`)
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&graph)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Cannot open graph file"}`)
	}

	dataFile.Close()


	fmt.Println("EDGES:", edgeFrom, edgeTo)


	best, err := graph.Shortest(int(edgeFrom), int(edgeTo))


	testWays := best.Path
	var testWayPairs []gravelmap.Way
	var prev = 0
	var routingData []gravelmap.RoutingLeg



	for i, testway := range testWays {
		if i == 0 {
			prev = testway
			continue
		}

		testWayPairs = append(testWayPairs, gravelmap.Way{EdgeFrom: int32(prev), EdgeTo: int32(testway)})

		prev = testway
	}

	wayFile, err := way.NewWayFileRead("_files")
	if err != nil {
		fmt.Fprintf(w, `{"message": "Cannot open way files"}`)
	}


	presentableWays, _ := wayFile.Read(testWayPairs)
	for _, pWay := range presentableWays {
		var latLngs []gravelmap.Point
		tmpLatLngs, _ := maps.DecodePolyline(pWay.Polyline)

		for _, latlng := range tmpLatLngs {
			latLngs = append(latLngs, gravelmap.Point{Lat: latlng.Lat, Lng: latlng.Lng})
		}

		routingLeg := gravelmap.RoutingLeg{
			Coordinates: latLngs,
			Length: 10,
			Paved: pWay.SurfaceType == gravelmap.WayTypeUnaved,
			Elevation: &gravelmap.RoutingLegElevation{Grade: float64(pWay.ElevationGrade), Start: 10, End: 20.0},
		}

		routingData = append(routingData, routingLeg)
	}

	json, _ := json.Marshal(routingData)
	fmt.Fprintf(w, string(json))
}

func getPointFromParamsNew(param string, r *http.Request) (*gravelmap.Point, error) {
	fromKeys, ok := r.URL.Query()[param]
	if !ok || len(fromKeys[0]) < 1 {
		return nil, errors.New("non existing param")
	}
	latLng := strings.Split(fromKeys[0], ",")

	lat, err := strconv.ParseFloat(latLng[0], 64)
	if err != nil {
		return nil, err
	}

	lng, err := strconv.ParseFloat(latLng[1], 64)
	if err != nil {
		return nil, err
	}

	return &gravelmap.Point{Lat: lat, Lng: lng}, nil
}