package commands

import (
	"fmt"
	"github.com/spf13/cobra"

	"database/sql"
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

    var centerLonLat = new OpenLayers.LonLat(-0.1279688, 51.5077286)
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
	return &cobra.Command{
		Use:   "map-drawer",
		Short: "Draw a route on a map",
		Long:  "Draw a route on a map and spit out the html",
		Run: func(cmd *cobra.Command, args []string) {
			connStr := "user=tkontos dbname=routing password=1234 port=5434"
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				fmt.Println("errore")
			}

			var source string
			findSourceQuery := `SELECT
	  source,
	  ST_Distance(ST_GeogFromText(CONCAT('SRID=4326;POINT(', x1, ' ', y1,')')), ST_GeogFromText('SRID=4326;POINT(23.8110783 38.0030367)')) as distance
	FROM ways
	ORDER BY distance
	LIMIT 1;`

			rows, err := db.Query(findSourceQuery)
			for rows.Next() {
				var row Way
				if err := rows.Scan(&row.node, &row.distance); err != nil {
					fmt.Println(err)
				} else {
					source = row.node
				}
			}





			var destination string
			findDestinationQuery := `SELECT
	  source,
	  ST_Distance(ST_GeogFromText(CONCAT('SRID=4326;POINT(', x2, ' ', y2,')')), ST_GeogFromText('SRID=4326;POINT(23.8312455 37.9495728)')) as distance
	FROM ways
	ORDER BY distance
	LIMIT 1;`

			rows, err = db.Query(findDestinationQuery)
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


			pointsArray := ""
			rows, err = db.Query(query)

			for rows.Next() {
				var row Row
				if err := rows.Scan(&row.pointFrom, &row.pointTo, &row.points); err != nil {
					fmt.Println(err)
				} else {
					s := strings.TrimPrefix(row.points, "LINESTRING(")
					s = strings.TrimSuffix(s, ")")
					sa := strings.Split(s, ",")

					for _, point := range sa {
						pointsArray += "[" + strings.Replace(point, " ", ", ", 1) + "],\n"
					}
				}
			}

			fmt.Println(fmt.Sprintf(html, pointsArray))
		},
	}
}
