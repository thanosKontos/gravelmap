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
	OSMID int64
	node int64
	source int64
	target int64
	elevationCost float64
	reverseElevationCost float64
	points string
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
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s", DBUser, DBPass, DBName, DBPort)
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
func (r *PgRouting) Route(pointFrom, pointTo gravelmap.Point) ([]gravelmap.RoutingFeature, error) {
	source, err := r.findClosestWaySourceId(pointFrom)
	if err != nil {
		return nil, err
	}

	destination, err := r.findClosestWaySourceId(pointTo)
	if err != nil {
		return nil, err
	}

	query := `SELECT
			osm_id, node, source, target, elevation_cost, reverse_elevation_cost, ST_AsText(the_geom) as points
		FROM pgr_dijkstra(
			'SELECT gid as id,
			source,
			target, 
			((CASE WHEN tag_id IN (101,201,202) THEN 1
				WHEN tag_id = 102 THEN 2
				WHEN tag_id = 203 THEN 3
				WHEN tag_id IN (103,104,105,106) THEN 5
            END)*elevation_cost*length_m) as cost,
	    	((CASE WHEN tag_id IN (101,201,202) THEN 1
				WHEN tag_id = 102 THEN 2
				WHEN tag_id = 203 THEN 3
				WHEN tag_id IN (103,104,105,106) THEN 5
            END)*reverse_elevation_cost*length_m) as reverse_cost
			FROM ways',
			%s,
			%s,
			TRUE
		) d INNER JOIN ways w ON d.edge = w.gid;`
	query = fmt.Sprintf(query, source, destination)

	r.logger.Debug(query)

	features := make([]gravelmap.RoutingFeature, 0)

	rows, err := r.routingClient.Query(query)
	for rows.Next() {
		coordinates := make([]gravelmap.Point, 0)
		feature := gravelmap.RoutingFeature{}

		var row routeRow
		if err := rows.Scan(&row.OSMID, &row.node, &row.source, &row.target, &row.elevationCost, &row.reverseElevationCost, &row.points); err != nil {
			return nil, err
		} else {
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
		}

		var elevationCost float64
		if row.node == row.source {
			elevationCost = row.elevationCost
		} else {
			elevationCost = row.reverseElevationCost
		}

		feature = gravelmap.RoutingFeature{
			Type: "LINESTRING",
			Coordinates: coordinates,
			Options: struct{
				OSMID int64
				ElevationCost float64
			}{
				row.OSMID,
				elevationCost,
			},
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
