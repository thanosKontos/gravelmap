package osm_grade

import (
	"database/sql"
	"fmt"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/elevation/srtm_ascii"
	"strconv"
	"strings"
	"sync"
)

const batchSize = 500

type geomRow struct {
	id     int64
	length float64
	geom   string
}

var wg sync.WaitGroup

type srtmOsmGrader struct {
	client    *sql.DB
	eleGrader *srtm_ascii.ElevationGrader
	logger     gravelmap.Logger
	OSMIDs     []string
}

// NewSRTM initialize and return an new SRTM object.
func NewSrtmOsmGrader(
	DBUser,
	DBPass,
	DBName,
	DBPort string,
	eleGrader *srtm_ascii.ElevationGrader,
	logger gravelmap.Logger,
) (*srtmOsmGrader, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s", DBUser, DBPass, DBName, DBPort)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &srtmOsmGrader{
		client:    DB,
		eleGrader: eleGrader,
		logger:    logger,
	}, nil
}

func (og *srtmOsmGrader) SetOSMIDs(OSMIDs []string) {
	og.OSMIDs = OSMIDs
}

func (og *srtmOsmGrader) GradeWays() error {
	og.client.Exec(`ALTER TABLE ways ADD COLUMN elevation_cost double precision DEFAULT 1;`)
	og.client.Exec(`ALTER TABLE ways ADD COLUMN reverse_elevation_cost double precision DEFAULT 1;`)
	og.client.Exec(`ALTER TABLE ways ADD COLUMN grade double precision DEFAULT NULL;`)
	og.client.Exec(`ALTER TABLE ways ADD COLUMN start_end_elevation text DEFAULT NULL;`)
	og.client.Exec(`ALTER TABLE ways ADD COLUMN way_graded boolean DEFAULT false;`)
	og.client.Exec(`CREATE INDEX idx_way_graded ON ways(way_graded);`)

	og.logger.Info("Added columns")

	rowsCount := 0
	countSQL := `SELECT COUNT(gid) FROM ways WHERE way_graded = false`
	waysSQL := `SELECT %s
		FROM ways
		WHERE way_graded = false
 		ORDER BY gid
		LIMIT %d`

	if len(og.OSMIDs) > 0 {
		countSQL = fmt.Sprintf("SELECT COUNT(gid) FROM ways WHERE osm_id IN (%s)", strings.Join(og.OSMIDs, ","))
		waysSQL = "SELECT %s FROM ways " + fmt.Sprintf("WHERE osm_id IN (%s) ORDER BY gid", strings.Join(og.OSMIDs, ",")) + " LIMIT %d"
	}

	row := og.client.QueryRow(countSQL)
	row.Scan(&rowsCount)

	for offset := 0; offset < rowsCount; offset = offset + batchSize {
		selectSQL := fmt.Sprintf(waysSQL, "gid, length_m, st_astext( ST_Transform( the_geom, 4326))", batchSize)
		onlyGidSelectSQL := fmt.Sprintf(waysSQL, "gid", batchSize)
		updateWayGradedSQL := "UPDATE ways SET way_graded = true WHERE gid IN (" + onlyGidSelectSQL + ")"
		rows, _ := og.client.Query(selectSQL)
		_, err := og.client.Exec(updateWayGradedSQL)
		if err != nil {
			og.logger.Warning(fmt.Sprintf("Error exec query `%s`: %s", updateWayGradedSQL, err))
		}

		wg.Add(1)
		go og.gradeBatch(rows)
	}

	wg.Wait()

	return nil
}

func (og *srtmOsmGrader) gradeBatch(rows *sql.Rows) {
	defer wg.Done()

	for rows.Next() {
		var row geomRow
		err := rows.Scan(&row.id, &row.length, &row.geom)

		if err == nil {
			og.gradeWay(row)
		}
	}

	og.logger.Info("batch finished")
}

func (og *srtmOsmGrader) gradeWay(row geomRow) {
	points, _ := geomToPoints(row.geom)
	wayElevation, err := og.eleGrader.Grade(points, row.length)
	updateElevationSQL := ""
	if err == nil {
		updateElevationSQL = fmt.Sprintf(
			`
UPDATE ways 
SET elevation_cost = %f,
reverse_elevation_cost = %f,
grade = %f,
start_end_elevation = '%f,%f'
WHERE gid = %d;`,
			gradeToCost(wayElevation.Grade),
			gradeToCost(-1*wayElevation.Grade),
			wayElevation.Grade,
			wayElevation.Start,
			wayElevation.End,
			row.id,
		)

		_, err = og.client.Exec(updateElevationSQL)
		if err != nil {
			og.logger.Warning(err)
		}
	} else {
		og.logger.Warning(err)
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
	switch {
	case grade <= -1:
		return 0.7
	case grade <= 0:
		return 0.9
	case grade <= 1:
		return 1.3
	case grade <= 3:
		return 3
	case grade <= 5:
		return 5
	default:
		return 10
	}
}
