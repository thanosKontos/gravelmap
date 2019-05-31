package commands

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/elevation/srtm_ascii"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const batchSize = 80

type geomRow struct {
	id     int64
	length float64
	geom   string
}

var wg sync.WaitGroup

// createRoutingDataCommand defines the create route command.
func createGradeWaysCommand() *cobra.Command {
	var (
		OSMIDs string
	)

	createGradeWaysCmd := &cobra.Command{
		Use:   "grade-ways",
		Short: "grade ways in route database",
		Long:  "grade ways and fill up elevation cost in route database",
	}

	createGradeWaysCmd.Flags().StringVar(&OSMIDs, "osm_ids", "", "The osm input file.")

	createGradeWaysCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createGradeWaysCmdRun(OSMIDs)
	}

	return createGradeWaysCmd
}

// createGradeWaysCmdRun defines the command run actions.
func createGradeWaysCmdRun(OSMIDs string) error {
	eleFinder, _ := srtm_ascii.NewElevationFinder(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"))
	eleGrader, _ := srtm_ascii.NewElevationGrader(eleFinder)

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"))
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	DB.Exec(`ALTER TABLE ways ADD COLUMN elevation_cost double precision DEFAULT 1;`)
	DB.Exec(`ALTER TABLE ways ADD COLUMN reverse_elevation_cost double precision DEFAULT 1;`)
	DB.Exec(`ALTER TABLE ways ADD COLUMN grade double precision DEFAULT NULL;`)

	rowsCount := 0
	countSQL := `SELECT COUNT(gid) FROM ways WHERE elevation_cost = 1 OR reverse_elevation_cost = 1`
	waysSQL := `SELECT gid, length_m, st_astext( ST_Transform( the_geom, 4326))
		FROM ways
		WHERE elevation_cost = 1 OR reverse_elevation_cost = 1
 		ORDER BY gid
		LIMIT %d OFFSET %d;`

	if OSMIDs != "" {
		countSQL = fmt.Sprintf("SELECT COUNT(gid) FROM ways WHERE osm_id IN (%s)", OSMIDs)
		waysSQL = fmt.Sprintf("SELECT gid, length_m, st_astext( ST_Transform( the_geom, 4326)) FROM ways WHERE osm_id IN (%s) ORDER BY gid", OSMIDs) + " LIMIT %d OFFSET %d;"
	}

	row := DB.QueryRow(countSQL)
	row.Scan(&rowsCount)

	for offset := 0; offset < rowsCount; offset = offset + batchSize {
		rows, _ := DB.Query(fmt.Sprintf(waysSQL, batchSize, offset))
		for rows.Next() {
			var row geomRow
			if err := rows.Scan(&row.id, &row.length, &row.geom); err != nil {
				return err
			} else {
				wg.Add(1)
				go gradeBatch(row, eleGrader, DB)
			}
		}

		wg.Wait()
		log.Println("batch finished")
	}

	log.Println("Roads graded.")

	return nil
}

func gradeBatch(row geomRow, eleGrader *srtm_ascii.ElevationGrader, DB *sql.DB) {
	defer wg.Done()
	points, _ := geomToPoints(row.geom)
	grade, err := eleGrader.Grade(points, row.length)
	if err == nil {
		updateElevationSQL := fmt.Sprintf(
			`UPDATE ways SET elevation_cost = %f, reverse_elevation_cost = %f, grade = %f WHERE gid = %d;`,
			gradeToCost(grade),
			gradeToCost(-1*grade),
			grade,
			row.id,
		)
		_, err := DB.Exec(updateElevationSQL)
		log.Println(fmt.Sprintf("Grade: %.2f", grade), "%")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println(err)
	}
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
	if grade <= -1 {
		return 0.7
	}

	if grade <= 0 {
		return 0.9
	}

	if grade <= 1 {
		return 1.3
	}

	if grade <= 3 {
		return 3
	}

	if grade <= 5 {
		return 5
	}

	return 10
}
