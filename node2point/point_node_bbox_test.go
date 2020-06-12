package node2point

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/distance"
)

type buffer struct {
	bytes.Buffer
}

// Add a Close method to our buffer so that we satisfy io.ReadWriteCloser.
func (b *buffer) Close() error {
	return nil
}

type bboxMemStorer struct {
	buffer buffer
}

func (ms *bboxMemStorer) getPointBbox(pt gravelmap.Point) string {
	// Level 1 bbox
	n := math.Floor(pt.Lat)
	e := math.Floor(pt.Lng)

	// Level 2 bbox
	l2 := math.Floor((pt.Lat - n) * 10)

	return fmt.Sprintf("N%.0fE%.0f_%.0f", n, e, l2)
}

func (ms *bboxMemStorer) getBboxWriteCloser(bbox string) (io.WriteCloser, error) {
	return &ms.buffer, nil
}

func (ms *bboxMemStorer) getPointReadCloser(pt gravelmap.Point) (io.ReadCloser, error) {
	return &ms.buffer, nil
}

func TestWriteBatchAndFindClosestNode(t *testing.T) {
	memStorer := &bboxMemStorer{}

	nodePointStore := &nodePointStore{
		nodePointBboxStorer: memStorer,
	}

	nodePointRead := &nodePointRead{
		distanceCalc:        distance.NewHaversine(),
		nodePointBboxStorer: memStorer,
	}

	nodePoints := []gravelmap.NodePoint{
		{10, gravelmap.Point{10.2, 11.3}},
		{11, gravelmap.Point{10.5, 11.3}},
		{12, gravelmap.Point{10.52, 11.4}},
	}

	err := nodePointStore.BatchStore(nodePoints)
	assert.Nil(t, err)

	closestNode, err := nodePointRead.FindClosest(gravelmap.Point{10.42, 11.3})
	assert.Nil(t, err)
	assert.Equal(t, int32(11), closestNode)
}

func TestFindClosestNotFindingNode(t *testing.T) {
	memStorer := &bboxMemStorer{}

	nodePointRead := &nodePointRead{
		distanceCalc:        distance.NewHaversine(),
		nodePointBboxStorer: memStorer,
	}

	_, err := nodePointRead.FindClosest(gravelmap.Point{10.42, 11.3})
	assert.Equal(t, "no node found", err.Error())
}
