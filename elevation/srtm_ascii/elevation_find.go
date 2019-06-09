package srtm_ascii

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/thanosKontos/gravelmap"
)

type nearbyElevation struct {
	elevation        int64
	distance         float64
	reversedDistance float64
}

type elevationRow struct {
	lng       float64
	lat       float64
	elevation int64
	distance  float64
}

type SRTMElevationFinder struct {
	client *sql.DB
}

// NewSRTM initialize and return an new SRTM object.
func NewElevationFinder(DBUser, DBPass, DBName, DBPort string) (*SRTMElevationFinder, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s", DBUser, DBPass, DBName, DBPort)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &SRTMElevationFinder{
		client: DB,
	}, nil
}

func (s *SRTMElevationFinder) FindElevation(point gravelmap.Point) (float64, error) {
	ptElevationSQL := `SELECT lng, lat, elevation_m, ST_Distance('SRID=4326;POINT(%f %f)'::geometry, geom) as distance
FROM elevation
ORDER BY geom <-> 'SRID=4326;POINT(%f %f)'::geometry
LIMIT 5;`

	ptElevationSQL = fmt.Sprintf(
		ptElevationSQL,
		point.Lng,
		point.Lat,
		point.Lng,
		point.Lat,
	)

	rows, err := s.client.Query(ptElevationSQL)
	if err != nil {
		return 0.0, err
	}

	var nearbyElevations []nearbyElevation
	var overallNearbyDistance float64
	count := 0
	for rows.Next() {
		var row elevationRow
		if err := rows.Scan(&row.lng, &row.lat, &row.elevation, &row.distance); err != nil {
			return 0.0, err
		} else {
			if row.distance > 35 {
				continue
			}

			nearbyElevations = append(nearbyElevations, nearbyElevation{elevation: row.elevation, distance: row.distance})
			overallNearbyDistance += row.distance
			count++
		}
	}
	if count < 4 {
		return 0.0, errors.New("not enough rows to calculate elevation")
	}

	var overallReversedDistance float64
	for i, nearbyElevation := range nearbyElevations {
		reversedDistance := overallNearbyDistance - nearbyElevation.distance
		nearbyElevations[i].reversedDistance = reversedDistance
		overallReversedDistance += reversedDistance
	}

	approximateElevation := 0.0
	for _, nearbyElevation := range nearbyElevations {
		approximateElevation += nearbyElevation.reversedDistance / overallReversedDistance * float64(nearbyElevation.elevation)
	}

	return approximateElevation, nil
}
