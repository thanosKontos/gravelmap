package commands

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"net/http"
	"strconv"
	"strings"
)

// createRoutingDataCommand defines the create routing command.
func createServerCommand() *cobra.Command {
	createServerCmd := &cobra.Command{
		Use:   "create-server",
		Short: "create a simple server to host a test routing website",
		Long:  "create a simple server to host a test routing website",
	}

	createServerCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createServerCmdRun()
	}

	return createServerCmd
}

// createRoutingDataCmdRun defines the command run actions.
func createServerCmdRun() error {
	http.HandleFunc("/route", func(w http.ResponseWriter, r *http.Request) {
		pointFrom, err := getPointFromParams("from", r)
		pointTo, err2 := getPointFromParams("to", r)
		if err != nil || err2 != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, `{"message": "Wrong arguments"}`)

			return
		}

		fmt.Println(pointFrom)
		fmt.Println(pointTo)
	})
	http.ListenAndServe(":8000", nil)

	return nil
}

func getPointFromParams(param string, r *http.Request) (*gravelmap.Point, error) {
	fromKeys, ok := r.URL.Query()[param]
	if !ok || len(fromKeys[0]) < 1 {
		return nil, errors.New("Non existing param")
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