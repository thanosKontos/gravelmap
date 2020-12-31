package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/graph"
	"github.com/thanosKontos/gravelmap/kml"
	"github.com/thanosKontos/gravelmap/node2point"
	"github.com/thanosKontos/gravelmap/path"
	"github.com/thanosKontos/gravelmap/route"
	"github.com/thanosKontos/gravelmap/routing"
	"github.com/thanosKontos/gravelmap/way"
)

// createWebServerCommand creates a web server for testing purposes.
func createWebServerCommand() *cobra.Command {
	createWebServerCmd := &cobra.Command{
		Use:   "create-web-server",
		Short: "create a simple server to host a test route website",
		Long:  "create a simple server to host a test route website",
	}

	createWebServerCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createWebServerCmdRun()
	}

	return createWebServerCmd
}

// createRoutingDataCmdRun defines the command run actions.
func createWebServerCmdRun() error {
	repo := graph.NewGobRepo("_files")

	mtbGraph, err := repo.Fetch("mtb")
	if err != nil {
		return err
	}

	hikeGraph, err := repo.Fetch("hike")
	if err != nil {
		return err
	}

	graphs := map[string]*graph.WeightedBidirectionalGraph{
		"mtb":  mtbGraph,
		"hike": hikeGraph,
	}

	http.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		routeHandler(w, r, graphs)
	})
	http.HandleFunc("/create-kml", func(w http.ResponseWriter, r *http.Request) {
		createKmlHandler(w, r, graphs)
	})

	http.ListenAndServe(":8000", nil)

	return nil
}

var errWrongArguments = errors.New("wrong arguments")
var errRouting = errors.New("error while routing")

func buildRoutingLegsFromRequestParams(r *http.Request, graphs map[string]*graph.WeightedBidirectionalGraph) ([]gravelmap.RoutingLeg, error) {
	pointFrom, err := getPointFromParams("from", r)
	pointTo, err2 := getPointFromParams("to", r)
	if err != nil || err2 != nil {
		return []gravelmap.RoutingLeg{}, errWrongArguments
	}

	routingModeParam, ok := r.URL.Query()["routing_mode"]
	if !ok || len(routingModeParam) != 1 || (routingModeParam[0] != "mtb" && routingModeParam[0] != "hike") {
		return []gravelmap.RoutingLeg{}, errWrongArguments
	}
	routingMode := routingModeParam[0]
	graph := graphs[routingMode]
	dijkstra := routing.NewDijsktraShortestPath(graph)

	edgeReader, err := way.NewWayFileRead("_files")
	if err != nil {
		return []gravelmap.RoutingLeg{}, errRouting
	}
	distanceCalc := distance.NewHaversine()
	edgeFinder := node2point.NewNodePointBboxFileRead("_files", distanceCalc)
	pathDecoder := path.NewGooglePolyline()
	router := route.NewGmRouter(edgeFinder, dijkstra, edgeReader, pathDecoder)
	routingData, err := router.Route(*pointFrom, *pointTo)
	if err != nil {
		err = errRouting
	}

	return routingData, err
}

func routeHandler(w http.ResponseWriter, r *http.Request, graphs map[string]*graph.WeightedBidirectionalGraph) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	routingData, err := buildRoutingLegsFromRequestParams(r, graphs)
	if err == errWrongArguments {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}
	if err == errRouting {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message": "Error routing"}`)

		return
	}

	json, _ := json.Marshal(routingData)
	fmt.Fprintf(w, string(json))
}

func createKmlHandler(w http.ResponseWriter, r *http.Request, graphs map[string]*graph.WeightedBidirectionalGraph) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	routingData, err := buildRoutingLegsFromRequestParams(r, graphs)
	if err == errWrongArguments {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}
	if err == errRouting {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message": "Error routing"}`)

		return
	}

	w.Header().Set("Content-Type", "application/vnd.google-earth.kml+xml")
	w.Header().Set("Content-Disposition", "attachment; filename=\"gravelmap.kml\"")

	kml := kml.NewKml()
	err = kml.Write(w, routingData)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message": "Error creating kml"}`)
	}
}

func getPointFromParams(param string, r *http.Request) (*gravelmap.Point, error) {
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
