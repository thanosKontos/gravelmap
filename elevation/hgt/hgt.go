package hgt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/thanosKontos/gravelmap"
)

const oneArcSecondRowColCount = 3601

var errorCannotGradeWay = errors.New("could not grade way")

type unzipper interface {
	unzip(zipFilename string) error
}

type downloader interface {
	download(dms string) error
}

type hgt struct {
	files          map[string]*os.File
	destinationDir string
	unzipper       unzipper
	downloader     downloader
	logger         gravelmap.Logger
}

// NewHgt instanciates a new HGT object
func NewHgt(destinationDir, nasaUsername, nasaPassword string, logger gravelmap.Logger) *hgt {
	unzipper := &unzip{destinationDir}
	downloader := &download{nasaUsername, nasaPassword, destinationDir}
	return &hgt{
		files:          make(map[string]*os.File),
		destinationDir: destinationDir,
		unzipper:       unzipper,
		downloader:     downloader,
		logger:         logger,
	}
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
		file, err := h.getFile(dms)
		if err != nil {
			return nil, err
		}

		latDiff := pt.Lat - math.Floor(pt.Lat)
		lngDiff := pt.Lng - math.Floor(pt.Lng)

		row := oneArcSecondRowColCount - int64(math.Round(latDiff*oneArcSecondRowColCount))
		col := int64(math.Round(lngDiff * oneArcSecondRowColCount))

		position := row*oneArcSecondRowColCount + col

		file.Seek(position*2, 0)

		data, err := readNextBytes(file, 2)
		if err != nil {
			return nil, err
		}
		buffer := bytes.NewBuffer(data)
		d := make([]byte, 2)

		err = binary.Read(buffer, binary.BigEndian, d)
		if err != nil {
			return nil, err
		}

		ele := int32(binary.BigEndian.Uint16(d))
		if ele > 60000 {
			h.logger.Debug("Could not grade (wrong elevation). Probably water, will use 0 instead")

			ele = 0
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
		Elevations: ptElevations,
		BidirectionalElevationInfo: gravelmap.BidirectionalElevationInfo{
			Normal:  gravelmap.ElevationInfo{Grade: grade, From: elevationStart, To: elevationEnd},
			Reverse: gravelmap.ElevationInfo{Grade: (-1) * grade, From: elevationEnd, To: elevationStart},
		},
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

func (h *hgt) downloadFile(dms string) error {
	h.logger.Info(fmt.Sprintf("Start downloading file: %s", dms))

	h.downloader.download(dms)
	h.unzipper.unzip(dms)

	h.logger.Info(fmt.Sprintf("Done downloading file: %s", dms))

	return nil
}

func (h *hgt) getFile(dms string) (*os.File, error) {
	if f, ok := h.files[dms]; ok {
		return f, nil
	}

	f, err := os.Open(fmt.Sprintf("%s/%s.hgt", h.destinationDir, dms))
	if err != nil {
		err = h.downloadFile(dms)
		if err != nil {
			return nil, err
		}

		f, err = os.Open(fmt.Sprintf("%s/%s.hgt", h.destinationDir, dms))
		if err != nil {
			return nil, err
		}
	}

	h.files[dms] = f

	return f, nil
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
	for _, f := range h.files {
		f.Close()
	}
}