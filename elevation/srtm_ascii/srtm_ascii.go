package srtm_ascii

import (
	"database/sql"
	"fmt"
	"github.com/thanosKontos/gravelmap"
)

// SRTM struct handles SRTM elevation.
type SRTM struct {
	filename string
	client   *sql.DB
	logger   gravelmap.Logger
}

// NewSRTM initialize and return an new SRTM object.
func NewSRTM(filename, DBUser, DBPass, DBName, DBPort string, logger gravelmap.Logger) (*SRTM, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=disable", DBUser, DBPass, DBName, DBPort)
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
