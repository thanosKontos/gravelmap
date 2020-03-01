package way

const (
	edgeStartFilename   = "edge_start.bin"
	edgeToFilename      = "edge_to_polylines_lookup.bin"
	polylinesFilename   = "polylines.bin"

	edgeStartRecordSize = 12

	// Each individual record has 3 int32s, 1 int64, 1 int8 and 1 float32
	edgeToIndividualRecordSize = 29
)

type edgeStartRecord struct {
	ConnectionsCnt int32
	NodeToOffset int64
}

type polylinePosition struct {
	length int32
	offset int64
}

type edgeToRecord struct {
	nodeTo int32
	distance int32
	wayType int8
	grade float32
	elevationStart int16
	elevationEnd int16
	polylinePosition
}
