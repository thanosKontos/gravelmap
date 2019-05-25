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

const batchSize = 100

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

	countSQL := `SELECT COUNT(gid) FROM ways WHERE elevation_cost = 1 OR reverse_elevation_cost = 1`
	row := DB.QueryRow(countSQL)
	wayCount := 0
	row.Scan(&wayCount)

	for offset := 0; offset < wayCount; offset = offset+batchSize {
		SQL := `SELECT gid, st_astext( ST_Transform( the_geom, 4326))
		FROM ways
		WHERE elevation_cost = 1 OR reverse_elevation_cost = 1
		LIMIT %d OFFSET %d;`

		rows, _ := DB.Query(fmt.Sprintf(SQL, batchSize, offset))
		for rows.Next() {
			var row geomRow
			if err := rows.Scan(&row.id, &row.geom); err != nil {
				return err
			} else {
				points, _ := geomToPoints(row.geom)
				grade, err := eleGrader.Grade(points)
				if err == nil {
					updateElevationSQL := fmt.Sprintf(`UPDATE ways SET elevation_cost = %f, reverse_elevation_cost = %f WHERE gid = %d;`, gradeToCost(grade), gradeToCost(-1*grade), row.id)
					DB.Exec(updateElevationSQL)
				}
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

func gradeToCost(grade float64) float64 {
	if grade <= 0 {
		return 0.7
	}

	if grade <= 5 {
		return 0.9
	}

	if grade <= 10 {
		return 1.5
	}

	if grade <= 20 {
		return 3
	}

	return 5
}
