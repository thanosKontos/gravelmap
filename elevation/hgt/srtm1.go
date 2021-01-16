package hgt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"os"

	"github.com/thanosKontos/gravelmap"
)

const oneArcSecondRowColCount = 3601

var errorWrongElevation = errors.New("wrong elevation")

type srtm1 struct {
	f *os.File
}

func NewStrm1(f *os.File) *srtm1 {
	return &srtm1{f: f}
}

func (s *srtm1) Get(pt gravelmap.Point) (int32, error) {
	latDiff := pt.Lat - math.Floor(pt.Lat)
	lngDiff := pt.Lng - math.Floor(pt.Lng)

	row := oneArcSecondRowColCount - int64(math.Round(latDiff*oneArcSecondRowColCount))
	col := int64(math.Round(lngDiff * oneArcSecondRowColCount))

	position := row*oneArcSecondRowColCount + col
	s.f.Seek(position*2, 0)
	data, err := readNextBytes(s.f, 2)
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
		return 0, errorWrongElevation
	}

	return ele, nil
}

func readNextBytes(file *os.File, number int) ([]byte, error) {
	bytes := make([]byte, number)

	_, err := file.Read(bytes)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func (s *srtm1) Close() error {
	return s.f.Close()
}
