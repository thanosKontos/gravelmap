package hgt

import (
	"errors"
	"fmt"
	"math"

	"github.com/thanosKontos/gravelmap"
)

var errorCannotGradeWay = errors.New("could not grade way")

// NewHgt instanciates a new HGT object with a file fetcher from US government
func NewHgt(efs gravelmap.ElevationFileStorer, logger gravelmap.Logger) *hgt {
	return &hgt{
		elevationFileStorage: efs,
		logger:               logger,
	}
}

type hgt struct {
	elevationFileStorage gravelmap.ElevationFileStorer
	logger               gravelmap.Logger
}

func (h *hgt) Get(points []gravelmap.Point, distance float64) (*gravelmap.WayElevation, error) {
	var ptElevations []int32
	var elevationStart, elevationEnd int16

	if distance <= 10 {
		h.logger.Debug("Could not grade (small distance)")
		return nil, errorCannotGradeWay
	}

	for i, pt := range points {
		dms := getDMSFromPoint(pt)
		elevationGetter, err := h.elevationFileStorage.Get(dms)
		if err != nil {
			return nil, err
		}
		ele, err := elevationGetter.Get(pt)

		if err != nil {
			if err == errorWrongElevation {
				h.logger.Debug("Could not grade (wrong elevation). Probably water, will use 0 instead")
			} else {
				return nil, err
			}
		}

		if i == 0 {
			elevationStart = int16(ele)
		}

		if i == len(points)-1 {
			elevationEnd = int16(ele)
		}

		ptElevations = append(ptElevations, ele)
	}

	grade := float32((elevationEnd-elevationStart)*100) / float32(distance)
	return &gravelmap.WayElevation{
		Elevations:    ptElevations,
		ElevationInfo: gravelmap.ElevationInfo{Grade: grade, From: elevationStart, To: elevationEnd},
	}, nil
}

// getDMSFromPoint extract the DMS format (e.g. N09E011) from point
func getDMSFromPoint(pt gravelmap.Point) string {
	latPfx := "N"
	if pt.Lat < 0 {
		latPfx = "S"
	}

	lngPfx := "E"
	if pt.Lng < 0 {
		lngPfx = "W"
	}

	return fmt.Sprintf("%s%02d%s%03d", latPfx, int8(math.Floor(pt.Lat)), lngPfx, int8(math.Floor(pt.Lng)))
}
