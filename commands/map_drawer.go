package commands

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"math"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type Row struct {
	pointFrom string
	pointTo string
	points string
}

type Way struct {
	node string
	distance string
}

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

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s port=%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPASS"),
		os.Getenv("DBPORT"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	findSrcSql := `SELECT
		  source,
		  ST_Distance(ST_GeogFromText(CONCAT('SRID=4326;POINT(', x1, ' ', y1,')')), ST_GeogFromText('SRID=4326;POINT(%s %s)')) as distance
		FROM ways
		ORDER BY distance
		LIMIT 1;`
	rows, err := db.Query(fmt.Sprintf(findSrcSql, from[0], from[1]))

	var source string
	for rows.Next() {
		var row Way
		if err := rows.Scan(&row.node, &row.distance); err != nil {
			fmt.Println(err)
		} else {
			source = row.node
		}
	}




	findDstnSql := `SELECT
	  source,
	  ST_Distance(ST_GeogFromText(CONCAT('SRID=4326;POINT(', x2, ' ', y2,')')), ST_GeogFromText('SRID=4326;POINT(%s %s)')) as distance
	FROM ways
	ORDER BY distance
	LIMIT 1;`
	rows, err = db.Query(fmt.Sprintf(findDstnSql, to[0], to[1]))

	var destination string
	for rows.Next() {
		var row Way
		if err := rows.Scan(&row.node, &row.distance); err != nil {
			fmt.Println(err)
		} else {
			destination = row.node
		}
	}



	query := `SELECT
		CONCAT(y1,',', x1) as point_from,
		CONCAT(y2, ',', x2) as point_to,
		ST_AsText(the_geom) as points
	FROM pgr_dijkstra(
		'SELECT gid as id, source, target, cost, reverse_cost FROM ways',
		%s,
		%s,
		FALSE
	) d INNER JOIN ways w ON d.edge = w.gid;`
	query = fmt.Sprintf(query, source, destination)

	pointsJSArray := ""
	rows, err = db.Query(query)
	for rows.Next() {
		var row Row
		if err := rows.Scan(&row.pointFrom, &row.pointTo, &row.points); err != nil {
			fmt.Println(err)
		} else {
			s := strings.TrimPrefix(row.points, "LINESTRING(")
			s = strings.TrimSuffix(s, ")")
			points := strings.Split(s, ",")
			rdPointsCnt := len(points)
			count := 0
			for _, point := range points {
				printPoint := false

				if rdPointsCnt >= 30 {
					if math.Mod(float64(count), 8) == 0 {
						printPoint = true
					}
				} else if rdPointsCnt >= 20 {
					if math.Mod(float64(count), 5) == 0 {
						printPoint = true
					}
				} else if rdPointsCnt >= 10 {
					if math.Mod(float64(count), 3) == 0 {
						printPoint = true
					}
				} else {
					printPoint = true
				}

				if printPoint {
					pointsJSArray += "[" + strings.Replace(point, " ", ", ", 1) + "],\n"
				}
				count++
			}
		}
	}

	fmt.Println(fmt.Sprintf(html, from[0], from[1], pointsJSArray))

	return nil
}

