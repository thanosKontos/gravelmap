package edge

import "github.com/thanosKontos/gravelmap"

const bBoxDir = "edge_bbox"

// bboxEdgeRecord is the data written in the bbox binary files
type bboxEdgeRecord struct {
	Pt gravelmap.Point
	EdgeID int32
}
