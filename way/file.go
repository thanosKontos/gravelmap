package way

const (
	edgeStartFilename = "edge_start.bin"
	edgeToFilename    = "edge_to_polylines_lookup.bin"
	polylinesFilename = "polylines.bin"

	edgeStartRecordSize = 9

	// Each individual record has 3 int32s, 2 int64, 2 int16 and 1 int8
	edgeToIndividualRecordSize = 33
)

type edgeStartRecord struct {
	ConnectionsCnt int8
	NodeToOffset   int64
}

type polylinePosition struct {
	length int32
	offset int64
}

type edgeToRecord struct {
	nodeTo   int32
	distance int32
	wayType  int8
	elevFrom int16
	elevTo   int16
	polylinePosition

	// for debug reasons
	osmID    int64
}
