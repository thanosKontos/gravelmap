package srtm_ascii

import (
	"fmt"
	"github.com/thanosKontos/gravelmap"
)

type nearbyElevation struct {
	elevation int64
	distance float64
	reversedDistance float64
}

type elevationRow struct {
	lng     float64
	lat float64
	elevation int64
	distance float64
}

func (s *SRTM) FindElevation(point gravelmap.Point) (float64, error) {
	ptElevationSQL := `SELECT *,
    ST_Distance(ST_GeogFromText(CONCAT('SRID=4326;POINT(', lng, ' ', lat,')')), ST_GeogFromText('SRID=4326;POINT(%f %f)')) as distance
	FROM elevation
	WHERE lat >= %f
	AND lat <= %f
	AND lng >= %f
	AND lng <= %f
	ORDER BY distance
	LIMIT 5;`

	ptElevationSQL = fmt.Sprintf(
		ptElevationSQL,
		point.Lng,
		point.Lat,
		point.Lat - 0.001,
		point.Lat + 0.001,
		point.Lng - 0.001,
		point.Lng + 0.001,
	)

	rows, err := s.client.Query(ptElevationSQL)
	if err != nil {
		return 0.0, err
	}

	var nearbyElevations []nearbyElevation
	var overallNearbyDistance float64
	for rows.Next() {
		var row elevationRow
		if err := rows.Scan(&row.lng, &row.lat, &row.elevation, &row.distance); err != nil {
			return 0.0, err
		} else {
			if row.distance > 30 {
				continue
			}

			nearbyElevations = append(nearbyElevations, nearbyElevation{elevation: row.elevation, distance: row.distance})
			overallNearbyDistance += row.distance
		}
	}

	var overallReversedDistance float64
	for i, nearbyElevation := range nearbyElevations {
		reversedDistance := overallNearbyDistance - nearbyElevation.distance
		nearbyElevations[i].reversedDistance = reversedDistance
		overallReversedDistance += reversedDistance
	}

	approximateElevation := 0.0
	for _, nearbyElevation := range nearbyElevations {
		approximateElevation += nearbyElevation.reversedDistance/overallReversedDistance*float64(nearbyElevation.elevation)
	}

	return approximateElevation, nil
}
