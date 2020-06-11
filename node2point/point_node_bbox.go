package node2point

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/thanosKontos/gravelmap"
)

type nodePointBboxStorer interface {
	getPointBbox(pt gravelmap.Point) string
	getPointReadCloser(pt gravelmap.Point) (io.ReadCloser, error)
	getBboxWriteCloser(bbox string) (io.WriteCloser, error)
}

type bboxFileRead struct {
	distanceCalc        gravelmap.DistanceCalculator
	nodePointBboxStorer nodePointBboxStorer
}

func (fr *bboxFileRead) FindClosest(point gravelmap.Point) (int32, error) {
	rc, err := fr.nodePointBboxStorer.getPointReadCloser(point)
	if err != nil {
		return 0, err
	}
	defer rc.Close()

	var closestNode int32 = 0
	var closestNodeDistance int64 = 0
	for {
		nodePoint := gravelmap.NodePoint{}
		data := readNextBytes(rc, 20)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.BigEndian, &nodePoint)

		if closestNode == 0 {
			closestNode = nodePoint.NodeID
			closestNodeDistance = fr.distanceCalc.Calculate(nodePoint.Pt, point)
		} else {
			d := fr.distanceCalc.Calculate(nodePoint.Pt, point)
			if closestNodeDistance > d {
				closestNode = nodePoint.NodeID
				closestNodeDistance = d
			}
		}

		if nodePoint.NodeID == 0 {
			if closestNode == 0 {
				return 0, errors.New("no node found")
			}

			break
		}
	}

	return closestNode, nil
}

func readNextBytes(r io.Reader, number int) []byte {
	byteSeq := make([]byte, number)
	_, _ = r.Read(byteSeq)

	return byteSeq
}

type bboxFileStore struct {
	nodePointBboxStorer nodePointBboxStorer
}

func NewBBoxFileStore(storageDir string) *bboxFileStore {
	nodePointBboxStorer := NewNodePointBboxFileStore(storageDir)

	return &bboxFileStore{
		nodePointBboxStorer: nodePointBboxStorer,
	}
}

func (fs *bboxFileStore) BatchStore(ndPts []gravelmap.NodePoint) error {
	ndBatchFileMap := map[string][]gravelmap.NodePoint{}
	for _, gmNd := range ndPts {
		bbox := fs.nodePointBboxStorer.getPointBbox(gmNd.Pt)
		ndBatchFileMap[bbox] = append(ndBatchFileMap[bbox], gmNd)
	}

	for bbox, batch := range ndBatchFileMap {
		err := fs.writeBatch(bbox, batch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *bboxFileStore) writeBatch(bbox string, ndPts []gravelmap.NodePoint) error {
	wc, err := fs.nodePointBboxStorer.getBboxWriteCloser(bbox)
	defer wc.Close()
	if err != nil {
		return err
	}

	var recs []gravelmap.NodePoint
	for _, ndPt := range ndPts {
		recs = append(recs, gravelmap.NodePoint{NodeID: ndPt.NodeID, Pt: ndPt.Pt})
	}

	var buf bytes.Buffer
	err = binary.Write(&buf, binary.BigEndian, recs)
	if err != nil {
		return err
	}

	_, err = wc.Write(buf.Bytes())
	return err
}
