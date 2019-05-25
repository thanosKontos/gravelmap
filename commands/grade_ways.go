package commands

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/elevation/srtm_ascii"
	"log"
	"os"
	"strconv"
	"strings"
)

type geomRow struct {
	id   int64
	geom string
}

// createRoutingDataCommand defines the create route command.
func createGradeWaysCommand() *cobra.Command {
	createGradeWaysCmd := &cobra.Command{
		Use:   "grade-ways",
		Short: "grade ways in route database",
		Long:  "grade ways and fill up elevation cost in route database",
	}

	createGradeWaysCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createGradeWaysCmdRun()
	}

	return createGradeWaysCmd
}

// createGradeWaysCmdRun defines the command run actions.
func createGradeWaysCmdRun() error {
	eleFinder, _ := srtm_ascii.NewElevationFinder(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"))
	distanceFinder := distance.NewHaversine()
	eleGrader, _ := srtm_ascii.NewElevationGrader(eleFinder, distanceFinder)

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s", os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"))
	DB, _ := sql.Open("postgres", connStr)
	SQL := `SELECT gid, st_astext( ST_Transform( the_geom, 4326))
    FROM ways
    WHERE 
	LIMIT 5;`

	SQL = `SELECT
			gid, ST_AsText(the_geom) as points
	FROM pgr_dijkstra(
			'SELECT gid as id, source, target, cost, reverse_cost FROM ways',
			41808,
			17338,
			FALSE
	) d INNER JOIN ways w ON d.edge = w.gid;`

	rows, _ := DB.Query(SQL)
	for rows.Next() {
		var row geomRow
		if err := rows.Scan(&row.id, &row.geom); err != nil {
			return err
		} else {
			points, _ := geomToPoints(row.geom)
			grade, err := eleGrader.Grade(points)
			if err == nil {
				fmt.Println(grade)
			}
		}
	}

	log.Println("Roads graded.")

	return nil
}

func geomToPoints(geom string) ([]gravelmap.Point, error) {
	points := make([]gravelmap.Point, 0)
	s := strings.TrimPrefix(geom, "LINESTRING(")
	s = strings.TrimSuffix(s, ")")
	pointsStr := strings.Split(s, ",")

	for _, point := range pointsStr {
		pointsSl := strings.Split(point, " ")
		lng, err := strconv.ParseFloat(pointsSl[0], 64)
		if err != nil {
			return nil, err
		}

		lat, err := strconv.ParseFloat(pointsSl[1], 64)
		if err != nil {
			return nil, err
		}

		p := gravelmap.Point{Lat: lat, Lng: lng}
		points = append(points, p)
	}

	return points, nil
}
