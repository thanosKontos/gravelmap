package route

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/thanosKontos/gravelmap"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type routeRow struct {
	node              int64
	tagId             int64
	source            int64
	target            int64
	length            float64
	grade             sql.NullFloat64
	startEndElevation sql.NullString
	points            string
}

type wayRow struct {
	node     string
	distance string
}

type PgRouting struct {
	routingClient *sql.DB
	logger        gravelmap.Logger
}

// NewRouting initialize and return an new PgRouting object.
func NewPgRouting(DBUser, DBPass, DBName, DBPort string, logger gravelmap.Logger) (*PgRouting, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", DBUser, DBPass, DBName, DBPort)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &PgRouting{
		routingClient: DB,
		logger:        logger,
	}, nil
}

func (r *PgRouting) Close() error {
	return r.routingClient.Close()
}

// Route calculates the route between 2 points and gives a slice of trip legs as features.
func (r *PgRouting) Route(pointFrom, pointTo gravelmap.Point, mode gravelmap.RoutingMode) ([]gravelmap.RoutingLeg, error) {
	source, err := r.findClosestWaySourceId(pointFrom)
	if err != nil {
		return nil, err
	}

	destination, err := r.findClosestWaySourceId(pointTo)
	if err != nil {
		return nil, err
	}

	query := `SELECT
			node, tag_id, source, target, length_m, start_end_elevation, grade, ST_AsText(the_geom) as points
		FROM pgr_dijkstra(
			'SELECT gid as id,
			source,
			target, 
			%s
			FROM ways',
			%s,
			%s,
			TRUE
		) d INNER JOIN ways w ON d.edge = w.gid;`
	query = fmt.Sprintf(query, costSelectsFromRoutingMode(mode), source, destination)

	r.logger.Debug(query)

	features := make([]gravelmap.RoutingLeg, 0)

	rows, err := r.routingClient.Query(query)
	for rows.Next() {
		coordinates := make([]gravelmap.Point, 0)
		feature := gravelmap.RoutingLeg{}

		var row routeRow
		err := rows.Scan(&row.node, &row.tagId, &row.source, &row.target, &row.length, &row.startEndElevation, &row.grade, &row.points)

		if err != nil {
			return nil, err
		}

		coordinates = make([]gravelmap.Point, 0)
		s := strings.TrimPrefix(row.points, "LINESTRING(")
		s = strings.TrimSuffix(s, ")")
		points := strings.Split(s, ",")

		for _, point := range points {
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
			coordinates = append(coordinates, p)
		}

		var elevationGrade float64
		hasElev := false
		elevationStart := 0.0
		elevationEnd := 0.0
		splittedElev := make([]string, 0)

		if row.node == row.source {
			if !row.grade.Valid {
				elevationGrade = 1
			} else {
				elevationGrade = row.grade.Float64
				if row.startEndElevation.Valid {
					splittedElev = strings.Split(row.startEndElevation.String, ",")
					elevationStart, _ = strconv.ParseFloat(splittedElev[0], 64)
					elevationEnd, _ = strconv.ParseFloat(splittedElev[1], 64)
					hasElev = true
				}
			}
		} else {
			if !row.grade.Valid {
				elevationGrade = 1
			} else {
				elevationGrade = -1 * row.grade.Float64
				if row.startEndElevation.Valid {
					splittedElev = strings.Split(row.startEndElevation.String, ",")
					elevationStart, _ = strconv.ParseFloat(splittedElev[1], 64)
					elevationEnd, _ = strconv.ParseFloat(splittedElev[0], 64)
					hasElev = true
				}
			}
		}

		if hasElev {
			feature = gravelmap.RoutingLeg{
				Coordinates: coordinates,
				Length:      row.length,
				Paved:       isRoadPaved(row.tagId),
				Elevation: &gravelmap.RoutingLegElevation{
					Grade: elevationGrade,
					Start: elevationStart,
					End:   elevationEnd,
				},
			}
		} else {
			feature = gravelmap.RoutingLeg{
				Coordinates: coordinates,
				Length:      row.length,
				Paved:       isRoadPaved(row.tagId),
				Elevation:   nil,
			}
		}

		features = append(features, feature)
	}

	return features, nil
}

func (r *PgRouting) findClosestWaySourceId(point gravelmap.Point) (string, error) {
	findSrcSql := `SELECT
			  source,
			  ST_Distance(ST_GeogFromText(CONCAT('SRID=4326;POINT(', x1, ' ', y1,')')), ST_GeogFromText('SRID=4326;POINT(%f %f)')) as distance
			FROM ways
			ORDER BY distance
			LIMIT 1;`

	row := r.routingClient.QueryRow(fmt.Sprintf(findSrcSql, point.Lat, point.Lng))

	var way wayRow
	if err := row.Scan(&way.node, &way.distance); err != nil {
		return "", err
	}

	distance, err := strconv.ParseFloat(way.distance, 64)
	if err != nil {
		return "", err
	}

	if distance > gravelmap.MinRoutingDistance {
		return "", errors.New("cannot find a close point")
	}

	return way.node, nil
}

func isRoadPaved(wayTagId int64) bool {
	pavedTagIds := [4]int64{103, 104, 105, 106}
	for _, p := range pavedTagIds {
		if p == wayTagId {
			return true
		}
	}
	return false
}

func costSelectsFromRoutingMode(mode gravelmap.RoutingMode) string {
	normalSurfaceCostFactor := "CASE WHEN tag_id IN (101,201,202) THEN 1 WHEN tag_id = 102 THEN 2 WHEN tag_id = 203 THEN 3 WHEN tag_id IN (103,104,105,106) THEN 5 END"
	onlyPavedSurfaceCostFactor := "CASE WHEN tag_id IN (101,201,202) THEN 1 WHEN tag_id = 102 THEN 2 WHEN tag_id = 203 THEN 100 WHEN tag_id IN (103,104,105,106) THEN 200 END"

	if mode == gravelmap.Normal {
		return fmt.Sprintf("((%s)*elevation_cost*length_m) as cost, ((%s)*reverse_elevation_cost*length_m) as reverse_cost", normalSurfaceCostFactor, normalSurfaceCostFactor)
	}

	if mode == gravelmap.OnlyUnpavedAccountElevation {
		return fmt.Sprintf("((%s)*elevation_cost*length_m) as cost, ((%s)*reverse_elevation_cost*length_m) as reverse_cost", onlyPavedSurfaceCostFactor, onlyPavedSurfaceCostFactor)
	}

	if mode == gravelmap.OnlyUnpavedHardcore {
		return fmt.Sprintf("((%s)*length_m) as cost, ((%s)*length_m) as reverse_cost", onlyPavedSurfaceCostFactor, onlyPavedSurfaceCostFactor)
	}

	if mode == gravelmap.NoLengthCareNormal {
		return fmt.Sprintf("((%s)*elevation_cost*length_m) as cost, ((%s)*reverse_elevation_cost) as reverse_cost", normalSurfaceCostFactor, normalSurfaceCostFactor)
	}

	return fmt.Sprintf("(%s) as cost, (%s) as reverse_cost", onlyPavedSurfaceCostFactor, onlyPavedSurfaceCostFactor)
}
