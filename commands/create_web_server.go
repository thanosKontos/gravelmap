package commands

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"net/http"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/graph"
	"github.com/thanosKontos/gravelmap/kml"
	"github.com/thanosKontos/gravelmap/node2point"
	"github.com/thanosKontos/gravelmap/route"
	"github.com/thanosKontos/gravelmap/routing_algorithm/dijkstra"
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
	mtbGraph := graph.NewGraph()
	dataFile, err := os.Open("_files/graph_bicycle.gob")
	if err != nil {
		return err
	}

	dataDecoder := gob.NewDecoder(dataFile)
	err = dataDecoder.Decode(&mtbGraph)
	if err != nil {
		return err
	}
	dataFile.Close()

	graphs := map[string]*graph.Graph{
		"bicycle": mtbGraph,
	}

	http.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		routeHandler(w, r, graphs)
	})
	http.HandleFunc("/create-kml", func(w http.ResponseWriter, r *http.Request) {
		createKmlHandler(w, r, mtbGraph)
	})

	http.ListenAndServe(":8000", nil)

	return nil
}

func routeHandler(w http.ResponseWriter, r *http.Request, graphs map[string]*graph.Graph) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pointFrom, err := getPointFromParams("from", r)
	pointTo, err2 := getPointFromParams("to", r)
	if err != nil || err2 != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}

	routingModeParam, ok := r.URL.Query()["routing_mode"]
	if !ok || len(routingModeParam) != 1 || (routingModeParam[0] != "bicycle" && routingModeParam[0] != "foot") {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}
	routingMode := routingModeParam[0]

	distanceCalc := distance.NewHaversine()
	edgeFinder := node2point.NewNodePointBboxFileRead("_files", distanceCalc)

	edgeReader, err := way.NewWayFileRead("_files")
	if err != nil {
		fmt.Fprintf(w, `{"message": "Cannot open way files"}`)

		return
	}

	graph := graphs[routingMode]
	dijkstra := dijkstra.NewDijkstra(graph)

	router := route.NewGmRouter(edgeFinder, dijkstra, edgeReader)
	routingData, err := router.Route(*pointFrom, *pointTo)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{"message": "%s"}`, err)

		return
	}

	json, _ := json.Marshal(routingData)
	fmt.Fprintf(w, string(json))
}

func createKmlHandler(w http.ResponseWriter, r *http.Request, graph *graph.Graph) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pointFrom, err := getPointFromParams("from", r)
	pointTo, err2 := getPointFromParams("to", r)
	if err != nil || err2 != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

		return
	}

	distanceCalc := distance.NewHaversine()
	edgeFinder := node2point.NewNodePointBboxFileRead("_files", distanceCalc)

	edgeReader, err := way.NewWayFileRead("_files")
	if err != nil {
		fmt.Fprintf(w, `{"message": "Cannot open way files"}`)

		return
	}

	dijkstra := dijkstra.NewDijkstra(graph)
	router := route.NewGmRouter(edgeFinder, dijkstra, edgeReader)
	routingData, err := router.Route(*pointFrom, *pointTo)
	if err != nil {
		w.WriteHeader(500)
		json, _ := json.Marshal(err)
		fmt.Fprintf(w, string(json))

		return
	}

	kml := kml.NewKml()
	kmlString, err := kml.CreateFromRoute(routingData)
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
