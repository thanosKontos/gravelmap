package way

const (
	edgeStartFilename   = "edge_start.bin"
	edgeToFilename      = "edge_to_polylines_lookup.bin"
	polylinesFilename   = "polylines.bin"

	edgeStartRecordSize = 12
)

type edgeStartRecord struct {
	ConnectionsCnt int32
	NodeToOffset int64
}
