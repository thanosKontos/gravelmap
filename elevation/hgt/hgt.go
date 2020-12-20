package hgt

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/thanosKontos/gravelmap"
)

var errorCannotGradeWay = errors.New("could not grade way")

type unzipper interface {
	unzip(zipFilename string) error
}

type hgtFileGetter interface {
	getFile(dms string) (*os.File, error)
}

type hgt struct {
	dmsElevationGettersCache map[string]gravelmap.ElevationPointGetterCloser
	logger                   gravelmap.Logger
	fileGetter               hgtFileGetter
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
		elevationGetter, err := h.getElevationGetter(dms)
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

func readNextBytes(file *os.File, number int) ([]byte, error) {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func (h *hgt) getElevationGetter(dms string) (gravelmap.ElevationPointGetter, error) {
	if g, ok := h.dmsElevationGettersCache[dms]; ok {
		return g, nil
	}

	h.logger.Info(fmt.Sprintf("Getting file: %s", dms))
	f, err := h.fileGetter.getFile(dms)
	if err != nil {
		h.logger.Error(err)
		return nil, err
	}
	h.logger.Info("Done")

	g := NewStrm1(f)
	h.dmsElevationGettersCache[dms] = g

	return g, nil
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

func (h *hgt) Close() {
	for _, egc := range h.dmsElevationGettersCache {
		egc.Close()
	}
}
