package way

const (
	edgeStartFilename   = "edge_start.bin"
	edgeToFilename      = "edge_to_polylines_lookup.bin"
	polylinesFilename   = "polylines.bin"

	edgeStartRecordSize = 12

	// Each individual record has 2 int32s, 1 int64, 1 int8 and 1 float32
	edgeToIndividualRecordSize = 21
)

type edgeStartRecord struct {
	ConnectionsCnt int32
	NodeToOffset int64
}
