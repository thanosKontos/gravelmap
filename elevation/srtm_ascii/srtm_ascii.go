package srtm_ascii

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/thanosKontos/gravelmap"
	"os"
	"strconv"
	"strings"
)

type ascData struct {
	lng       float64
	lat       float64
	elevation int
}

type fileInfo struct {
	headerCnt int
	rowsCnt   int
	colsCnt   int
	latMin    float64
	latMax    float64
	lngMin    float64
	lngMax    float64
	step      float64
	noDataVal string
}

// SRTM struct handles SRTM elevation.
type SRTM struct {
	filename string
	client   *sql.DB
	logger   gravelmap.Logger
}

// NewSRTM initialize and return an new SRTM object.
func NewSRTM(filename, DBUser, DBPass, DBName, DBPort string, logger gravelmap.Logger) (*SRTM, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s", DBUser, DBPass, DBName, DBPort)
	DB, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return &SRTM{
		filename: filename,
		client:   DB,
		logger:   logger,
	}, nil
}

func (s *SRTM) Import() error {
	err := s.createElevationTable()
	if err != nil {
		return nil
	}

	return s.processFile()
}

func (s *SRTM) processFile() error {
	info, err := s.extractHeaderInfo()
	if err != nil {
		return err
	}

	err = s.importToDB(info)

	return err
}

func (s *SRTM) importToDB(info *fileInfo) error {
	file, err := os.Open(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	lineCnt := 0
	lat := info.latMax
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		eleData := make([]ascData, 0)
		lng := info.lngMin
		lineCnt++
		if lineCnt <= info.headerCnt {
			continue
		}

		lineTxt := scanner.Text()
		lngElevations := strings.Fields(lineTxt)

		for _, eleStr := range lngElevations {
			lng += info.step
			ele, _ := strconv.Atoi(eleStr)
			eleData = append(eleData, ascData{lng: lng, lat: lat, elevation: ele})
		}

		lat -= info.step

		s.insertElevations(eleData)
	}

	return nil
}

func (s *SRTM) insertElevations(eleData []ascData) error {
	values := make([]string, 0)
	for _, e := range eleData {
		values = append(values, fmt.Sprintf(`('%f', '%f', '%d')`, e.lng, e.lat, e.elevation))
	}

	insertSQL := fmt.Sprintf(`INSERT INTO elevation ("lng", "lat", "elevation_m") VALUES %s`, strings.Join(values, ", "))

	rst, err := s.client.Exec(insertSQL)
	if err != nil {
		s.logger.Error(err)
		return err
	}

	r, _ := rst.RowsAffected()
	s.logger.Info(fmt.Sprintf("%d rows affected", r))

	return nil
}

func (s *SRTM) extractHeaderInfo() (*fileInfo, error) {
	file, err := os.Open(s.filename)
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
		rowsCnt:   fDataRowCount,
		colsCnt:   fDataColCount,
		latMin:    fLatMin,
		latMax:    fLatMin + float64(fDataRowCount)*fstep,
		lngMin:    fLngMin,
		lngMax:    fLngMin + float64(fDataColCount)*fstep,
		step:      fstep,
		noDataVal: info["NODATA_value"],
	}

	return &fi, nil
}

func (s *SRTM) createElevationTable() error {
	_, err := s.client.Exec(`CREATE TABLE IF NOT EXISTS public.elevation
	(
		lng double precision,
		lat double precision,
		elevation_m integer,
		PRIMARY KEY(lat, lng)
	)`)

	return err
}
