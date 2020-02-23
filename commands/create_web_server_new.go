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
	//testWays := []int32{1,86123,135138,135133,121181,85173,5519,121174,116378,85694,86138,63143,4689,85131,121195,86120,85760,112247,63577,112242,112237,135110,56424,85141,135102,56428,135077,973,132006,135067,82937,698,85158,135054,132060,13339,132055,107433,112180,134875,85124,115145,39299,132050,96635,138762,138765,152794,112184,152793,152792,1690,86599,86594,86592,86591,86581,121805,2170,2173}
	//testWays := []int32{1,86123,135138,135133,121181}
	var testWayPairs []gravelmap.Way
	var prev int = 0
	for i, testway := range testWays {
		if i == 0 {
			prev = testway
			continue
		}

		testWayPairs = append(testWayPairs, gravelmap.Way{int32(prev), int32(testway)})

		prev = testway
	}

	wayFile := way.NewWayFileRead("_files")
	polylines, _ := wayFile.Read(testWayPairs)

	var latLngs []struct{
		Lat float64
		Lng float64
	}
	for _, pl := range polylines {
		tmpLatLngs, _ := maps.DecodePolyline(pl)

		for _, latlng := range tmpLatLngs {
			latLngs = append(latLngs, struct{
				Lat float64
				Lng float64
			}{Lat: latlng.Lat, Lng: latlng.Lng})
		}
	}

	routingData := []struct {
		Elevation interface{}
		Paved bool
		Coordinates []struct{
			Lat float64
			Lng float64
		}
	}{{
		nil, true, latLngs,
	}}

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