package hgt

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"

	"github.com/thanosKontos/gravelmap"
)

const nasa30mSrtmURL = "http://e4ftl01.cr.usgs.gov/MEASURES/SRTMGL1.003/2000.02.11/%s.SRTMGL1.hgt.zip"

type downloader interface {
	download(dms string) error
}

// NewNasaHgt instanciates a new HGT object with files coming from nasa srtm servers
func NewNasaHgt(destinationDir, username, password string, logger gravelmap.Logger) *hgt {
	fileGetter := &nasa30mFile{username, password, destinationDir}

	return &hgt{
		dmsElevationGettersCache: make(map[string]gravelmap.ElevationPointGetterCloser),
		logger:                   logger,
		fileGetter:               fileGetter,
	}
}

type nasa30mFile struct {
	username       string
	password       string
	destinationDir string
}

func (n *nasa30mFile) getFile(dms string) (*os.File, error) {
	f, err := os.Open(fmt.Sprintf("%s/%s.hgt", n.destinationDir, dms))
	if err == nil {
		return f, err
	}

	err = n.download(dms)
	if err != nil {
		return nil, err
	}
	err = n.unzip(dms)
	if err != nil {
		return nil, err
	}

	return os.Open(fmt.Sprintf("%s/%s.hgt", n.destinationDir, dms))
}

func (n *nasa30mFile) download(dms string) error {
	out, err := os.Create(fmt.Sprintf("/%s/%s.hgt.zip", n.destinationDir, dms))
	defer out.Close()
	if err != nil {
		return err
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar:           jar,
		CheckRedirect: n.redirectPolicyFunc,
	}
	req, err := http.NewRequest("GET", fmt.Sprintf(nasa30mSrtmURL, dms), nil)
	req.SetBasicAuth(n.username, n.password)
	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (n *nasa30mFile) redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.SetBasicAuth(n.username, n.password)
	return nil
}

func (n *nasa30mFile) unzip(dms string) error {
	zipFilename := fmt.Sprintf("/%s/%s.hgt.zip", n.destinationDir, dms)
	zipReader, err := zip.OpenReader(zipFilename)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Archive is supposed to have just a single hgt file
	if len(zipReader.Reader.File) != 1 {
		return errors.New("wrong hgt zip formats")
	}

	file := zipReader.Reader.File[0]
	if file.FileInfo().IsDir() {
		return errors.New("wrong hgt zip format")
	}

	zippedFile, err := file.Open()
	if err != nil {
		return err
	}
	defer zippedFile.Close()

	extractedFilePath := filepath.Join(
		n.destinationDir,
		file.Name,
	)

	outputFile, err := os.OpenFile(
		extractedFilePath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		file.Mode(),
	)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, zippedFile)
	if err != nil {
		return err
	}

	return nil
}
