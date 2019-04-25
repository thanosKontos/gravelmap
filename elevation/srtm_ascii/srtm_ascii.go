package srtm_ascii

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// SRTM struct handles SRTM elevation.
type SRTM struct {
}

// NewSRTM initialize and return an new SRTM object.
func NewSRTM() *SRTM {
	return &SRTM{}
}

type fileInfo struct {
	headerCnt int
	rowsCnt int
	colsCnt int
	latMin float64
	lngMin float64
	step float64
	noDataVal string
}

// Find finds the elevation for a specific coordinate in meters
func (SRTM) Find(lat, lng float64) (int64, error) {
	filename := findFileToQuery(lat, lng)
	fInfo, err := extractInfoFromFile(filename)

	rowΝο := fInfo.headerCnt + fInfo.rowsCnt - int(math.Round((lat - fInfo.latMin)/fInfo.step))
	colNo := int(math.Round((lng - fInfo.lngMin)/fInfo.step))

	cmd := exec.Command("awk", fmt.Sprintf("NR == %d {print $%d}", rowΝο, colNo), filename)
	out, err := cmd.CombinedOutput()

	if err != nil {
		return 0, err
	}

	eleStr := strings.Trim(string(out), "\n")

	if eleStr == fInfo.noDataVal {
		return 0, errors.New("no data available for this point")
	}

	ele, err := strconv.ParseInt(eleStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return ele, nil
}

func extractInfoFromFile(filename string) (*fileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineTxt := scanner.Text()
		words := strings.Fields(lineTxt)
		if len(words) != 2 {
			break
		}

		info[words[0]] = words[1]
	}

	fDataRowCount, err := strconv.Atoi(info["nrows"])
	if err != nil {
		return nil, err
	}

	fDataColCount, err := strconv.Atoi(info["ncols"])
	if err != nil {
		return nil, err
	}

	fLatMin, err := strconv.ParseFloat(info["yllcorner"], 64)
	if err != nil {
		return nil, err
	}

	fLngMin, err := strconv.ParseFloat(info["xllcorner"], 64)
	if err != nil {
		return nil, err
	}

	fstep, err := strconv.ParseFloat(info["cellsize"], 64)
	if err != nil {
		return nil, err
	}

	fi := fileInfo{
		headerCnt: len(info),
		rowsCnt: fDataRowCount,
		colsCnt: fDataColCount,
		latMin: fLatMin,
		lngMin: fLngMin,
		step: fstep,
		noDataVal: info["NODATA_value"],
	}

	return &fi, nil
}

func findFileToQuery(lat, lng float64) string {
	return "/home/tkontos/Downloads/lower_greece_elevation_data/ascii/srtm_41_05.asc"
}
