package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/routing/pgrouting"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var html = `<html><body>
  <div id="mapdiv"></div>
  <script src="http://www.openlayers.org/api/OpenLayers.js"></script>
  <script>
    map = new OpenLayers.Map("mapdiv");
    map.addLayer(new OpenLayers.Layer.OSM());

    var centerLonLat = new OpenLayers.LonLat(%s, %s)
          .transform(
            new OpenLayers.Projection("EPSG:4326"),
            map.getProjectionObject()
          );

    var markers = new OpenLayers.Layer.Markers("Markers");
    map.addLayer(markers);

    var points = [
        %s
    ];

    points.forEach(function(point) {
        var lonLat = new OpenLayers.LonLat(point[0], point[1])
          .transform(
            new OpenLayers.Projection("EPSG:4326"),
            map.getProjectionObject()
          );
        markers.addMarker(new OpenLayers.Marker(lonLat));
    });

    map.setCenter(centerLonLat, 16);
  </script>
</body></html>`

// mapDrawerCommand is a util to create a test HTML page in order to help with the tedious task of manual testing the router.
func mapDrawerCommand() *cobra.Command {
	var (
		pointFrom string
		pointTo   string
	)

	mapDrawerCmd := &cobra.Command{
		Use:   "map-drawer",
		Short: "Draw a route on a map",
		Long:  "Draw a route on a map and spit out the html",
	}

	mapDrawerCmd.Flags().StringVar(&pointFrom, "from", "", "Lat,lng to route from.")
	mapDrawerCmd.Flags().StringVar(&pointTo, "to", "", "Lat,lng to route to.")

	mapDrawerCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return mapDrawerCmdRun(pointFrom, pointTo)
	}

	return mapDrawerCmd
}

// mapDrawerCmdRun defines the command run actions.
func mapDrawerCmdRun(pointFrom, pointTo string) error {
	from := strings.Split(pointFrom, ",")
	to := strings.Split(pointTo, ",")

	pgRouting, err := pgrouting.NewPgRouting(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"))
	if err != nil {
		log.Fatal(err)
	}

	latFrom, err := strconv.ParseFloat(from[0], 64)
	if err != nil {
		log.Fatal(err)
	}

	lngFrom, err := strconv.ParseFloat(from[1], 64)
	if err != nil {
		log.Fatal(err)
	}

	latTo, err := strconv.ParseFloat(to[0], 64)
	if err != nil {
		log.Fatal(err)
	}

	lngTo, err := strconv.ParseFloat(to[1], 64)
	if err != nil {
		log.Fatal(err)
	}

	tripLegs, err := pgRouting.Route(
		gravelmap.Point{Lat: latFrom, Lng: lngFrom},
		gravelmap.Point{Lat: latTo, Lng: lngTo},
	)
	if err != nil {
		log.Fatal(err)
	}

	pointsJSArray := ""
	for _, leg := range tripLegs {
		for _, point := range leg {
			pointsJSArray += fmt.Sprintf("[%f,%f],\n", point.Lng, point.Lat)
		}
	}

	fmt.Println(fmt.Sprintf(html, from[0], from[1], pointsJSArray))

	return nil
}
