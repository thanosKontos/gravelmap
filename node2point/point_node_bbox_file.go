package node2point

import (
	"fmt"
	"io"
	"math"
	"os"

	"github.com/thanosKontos/gravelmap"
)

const bBoxDir = "edge_bbox"

type nodePointBboxFile struct {
	storageDir string
}

// NewNodePointBboxFileRead instantiate a new nodePointBboxFile object for reading
func NewNodePointBboxFileRead(storageDir string, dc gravelmap.DistanceCalculator) *nodePointRead {
	nodePointBboxStorer := &nodePointBboxFile{
		storageDir: storageDir,
	}
	return &nodePointRead{
		distanceCalc:        dc,
		nodePointBboxStorer: nodePointBboxStorer,
	}
}

func NewNodePointBboxFileStore(storageDir string) *nodePointStore {
	os.RemoveAll(fmt.Sprintf("%s/%s", storageDir, bBoxDir))
	os.Mkdir(fmt.Sprintf("%s/%s", storageDir, bBoxDir), 0777)
	nodePointBboxStorer := &nodePointBboxFile{
		storageDir: storageDir,
	}

	return &nodePointStore{
		nodePointBboxStorer: nodePointBboxStorer,
	}
}

func (bbf *nodePointBboxFile) getPointBbox(pt gravelmap.Point) string {
	// Level 1 bbox
	n := math.Floor(pt.Lat)
	e := math.Floor(pt.Lng)

	// Level 2 bbox
	l2 := math.Floor((pt.Lat - n) * 10)

	return fmt.Sprintf("N%.0fE%.0f_%.0f", n, e, l2)
}

func (bbf *nodePointBboxFile) getBboxWriteCloser(bbox string) (io.WriteCloser, error) {
	return os.OpenFile(
		fmt.Sprintf("%s/%s/%s.bin", bbf.storageDir, bBoxDir, bbox),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0777,
	)
}

func (bbf *nodePointBboxFile) getPointReadCloser(pt gravelmap.Point) (io.ReadCloser, error) {
	fname := bbf.getPointBbox(pt)
	return os.Open(fmt.Sprintf("%s/%s/%s.bin", bbf.storageDir, bBoxDir, fname))
}
