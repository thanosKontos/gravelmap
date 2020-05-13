package elevation

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"os"
	"os/exec"

	"github.com/thanosKontos/gravelmap"
)

const oneArcSecondRowColCount = 3601

var errorCannotGradeWay = errors.New("could not grade way")

type hgt struct {
	files          map[string]*os.File
	destinationDir string
	nasaUsername   string
	nasaPassword   string
	distanceCalc   gravelmap.DistanceCalculator
	logger         gravelmap.Logger
}

func NewHgt(
	destinationDir,
	nasaUsername,
	nasaPassword string,
	distanceCalc gravelmap.DistanceCalculator,
	logger gravelmap.Logger,
	) *hgt {
	return &hgt{
		files:          make(map[string]*os.File),
		destinationDir: destinationDir,
		nasaUsername:   nasaUsername,
		nasaPassword:   nasaPassword,
		distanceCalc:   distanceCalc,
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
		ele, err := h.getPointElevation(pt)
		if err != nil {
			return nil, err
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

type closebyEle struct {
	distance int64
	ele int32
	weight float64
}

func (h *hgt) getPointElevation(pt gravelmap.Point) (int32, error) {
	dms := getDMSFromPoint(pt)
	f, err := h.getFile(dms)
	if err != nil {
		return 0, err
	}

	baseLat := math.Floor(pt.Lat)
	baseLng := math.Floor(pt.Lng)

	latDiff := pt.Lat - baseLat
	lngDiff := pt.Lng - baseLng

	rowMax := oneArcSecondRowColCount - int64(math.Floor(latDiff*oneArcSecondRowColCount))
	rowMin := rowMax - 1
	if rowMax == 0 {
		rowMin = 0
	}

	colMin := int64(math.Floor(lngDiff * oneArcSecondRowColCount))
	colMax := colMin + 1
	if colMax == oneArcSecondRowColCount {
		colMax = oneArcSecondRowColCount
	}

	minLat := baseLat + (float64(3601-rowMax) * (1.0/3601))
	maxLat := baseLat + (float64(3601-rowMin) * (1.0/3601))
	minLng := baseLng + (1.0/3601)*float64(colMin)
	maxLng := baseLng + (1.0/3601)*float64(colMax)

	var closebyElevations []closebyEle
	topLeftEle, err := h.getEleFileRecord(f, rowMax, colMax)
	if err == nil {
		closebyElevations = append(closebyElevations, closebyEle{h.distanceCalc.Calculate(pt, gravelmap.Point{Lat: maxLat, Lng: minLng}), topLeftEle, 0.0})
	}

	topRightEle, err := h.getEleFileRecord(f, rowMax, colMin)
	if err == nil {
		closebyElevations = append(closebyElevations, closebyEle{h.distanceCalc.Calculate(pt, gravelmap.Point{Lat: maxLat, Lng: maxLng}), topRightEle, 0.0})
	}

	bottomLeftEle, err := h.getEleFileRecord(f, rowMin, colMax)
	if err == nil {
		closebyElevations = append(closebyElevations, closebyEle{h.distanceCalc.Calculate(pt, gravelmap.Point{Lat: minLat, Lng: minLng}), bottomLeftEle, 0.0})
	}

	bottomRightEle, err := h.getEleFileRecord(f, rowMin, colMin)
	if err == nil {
		closebyElevations = append(closebyElevations, closebyEle{h.distanceCalc.Calculate(pt, gravelmap.Point{Lat: minLat, Lng: maxLng}), bottomRightEle, 0.0})
	}

	largestDist := int64(0)
	for _, closebyEle := range closebyElevations {
		if closebyEle.distance > largestDist {
			largestDist = closebyEle.distance
		}
	}

	totalWeights := 0.0
	for i, closebyEle := range closebyElevations {
		closebyElevations[i].weight = float64(largestDist-closebyEle.distance)*100.0/float64(largestDist)
		totalWeights += closebyElevations[i].weight
	}

	elev := 0.0
	for _, closebyEle := range closebyElevations {
		elev += float64(closebyEle.ele)*closebyEle.weight/totalWeights
	}

	return int32(elev), nil
}

func (h *hgt) getEleFileRecord(f *os.File, row, col int64) (int32, error) {
	pos := row*oneArcSecondRowColCount + col
	f.Seek(pos*2, 0)

	data, err := readNextBytes(f, 2)
	if err != nil {
		return 0, err
	}
	buffer := bytes.NewBuffer(data)
	d := make([]byte, 2)

	err = binary.Read(buffer, binary.BigEndian, d)
	if err != nil {
		return 0, err
	}

	ele := int32(binary.BigEndian.Uint16(d))
	if ele > 60000 {
		h.logger.Debug("Could not grade (wrong elevation). Probably water, will use 0 instead")

		ele = 0
	}

	return ele, nil
}

func readNextBytes(file *os.File, number int) ([]byte, error) {
	bts := make([]byte, number)
	_, err := file.Read(bts)
	if err != nil {
		return []byte{}, err
	}

	return bts, nil
}

// TODO: replace the real wget and tar commands with net/http and archive/zip in order to be testable and less dependent
func (h *hgt) downloadFile(dms string) error {
	h.logger.Info(fmt.Sprintf("Start downloading file: %s", dms))

	url := fmt.Sprintf("http://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL1.003/2000.02.11/%s.SRTMGL1.hgt.zip", dms)
	out, err := exec.Command("wget", url, fmt.Sprintf("--http-user=%s", h.nasaUsername), fmt.Sprintf("--http-password=%s", h.nasaPassword), "-P", h.destinationDir).Output()
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
