package elevation

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"log"
	"math"

	"os"
	"os/exec"

	"github.com/thanosKontos/gravelmap"
)

const oneArcSecondRowColCount = 3601

var errorCannotGradeWay = errors.New("could not grade way")

type hgt struct {
	files map[string]*os.File
	destinationDir string
	logger gravelmap.Logger
}

func NewHgt(destinationDir string, logger gravelmap.Logger) *hgt {
	return &hgt{
		files: make(map[string]*os.File),
		destinationDir: destinationDir,
		logger: logger,
	}
}

func (h *hgt) Get(points []gravelmap.Point, distance float64) (*gravelmap.WayElevation, error) {
	var ptElevations []int32
	var prevEle int32
	var incline int32 = 0

	if distance <= 10 {
		h.logger.Warning("Could not grade (small distance)")
		return nil, errorCannotGradeWay
	}

	for i, pt := range points {
		dms := getDMSFromPoint(pt)
		file, err := h.getFile(dms)
		if err != nil {
			return nil, err
		}

		latDiff := pt.Lat-math.Floor(pt.Lat)
		lngDiff := pt.Lng-math.Floor(pt.Lng)

		row := oneArcSecondRowColCount-int64(math.Round(latDiff * oneArcSecondRowColCount))
		col := int64(math.Round(lngDiff * oneArcSecondRowColCount))

		position := row*oneArcSecondRowColCount + col

		file.Seek(position*2, 0)

		data := readNextBytes(file, 2)
		buffer := bytes.NewBuffer(data)
		d := make([]byte, 2)

		err = binary.Read(buffer, binary.BigEndian, d)
		if err != nil {
			return nil, err
		}

		ele := int32(binary.BigEndian.Uint16(d))
		if ele > 60000 {
			h.logger.Warning("Could not grade (void elevation)")
			return nil, errorCannotGradeWay
		}

		if i != 0 {
			incline += ele - prevEle
		}

		ptElevations = append(ptElevations, ele)
		prevEle = ele
	}

	return &gravelmap.WayElevation{Elevations: ptElevations, Incline: incline, Grade: float64(incline*100)/distance}, nil
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

// TODO: replace the real wget and tar commands with net/http and archive/zip in order to be testable and less dependent
func (h *hgt) downloadFile(dms string) error {
	h.logger.Debug(fmt.Sprintf("Start downloading file: %s", dms))

	url := fmt.Sprintf("http://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL1.003/2000.02.11/%s.SRTMGL1.hgt.zip", dms)
	out, err := exec.Command("wget", url, "--http-user=tkontos", "--http-password=1234", "-P", h.destinationDir).Output()
	if err != nil {
		h.logger.Error("wget error")

		return err
	}

	h.logger.Debug(string(out))



	zipFile := fmt.Sprintf("/%s/%s.SRTMGL1.hgt.zip", h.destinationDir, dms)
	out, err = exec.Command("tar", "xvf", zipFile, "-C", h.destinationDir).Output()
	if err != nil {
		h.logger.Error("untar error")

		return err
	}

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

func (h *hgt) Close () {
	for _, f := range h.files {
		f.Close()
	}
}