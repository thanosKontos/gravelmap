package hgt

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type unzip struct {
	destinationDir string
}

func (uz *unzip) unzip(dms string) error {
	zipFilename := fmt.Sprintf("/%s/%s.hgt.zip", uz.destinationDir, dms)
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
		uz.destinationDir,
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
