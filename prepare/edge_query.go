package prepare

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/go-memdb"
	"github.com/qedus/osmpbf"
)

const edgeTable = "gravelmap_osm_node_count"

type edgeQuery struct {
	osmFilename  string
	db *memdb.MemDB
}

type osmNodeCount struct {
	ID    int64
	Count int
}

func NewEdgeQuery(osmFilename string) *edgeQuery {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			edgeTable: {
				Name: edgeTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
					"count": {
						Name:    "count",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Count"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return &edgeQuery{
		osmFilename:  osmFilename,
		db: db,
	}
}

func (e *edgeQuery) Prepare () {
	f, err := os.Open(e.osmFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)

	// use more memory from the start, it is faster
	d.SetBufferSize(osmpbf.MaxBlobSize)

	// start decoding with several goroutines, it is faster
	err = d.Start(runtime.GOMAXPROCS(-1))
	if err != nil {
		log.Fatal(err)
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				wtTxn := e.db.Txn(true)

				for _, nodeID := range v.NodeIDs {
					rdTxn := e.db.Txn(false)
					raw, err := rdTxn.First(edgeTable, "id", nodeID)
					rdTxn.Abort()

					if err == nil && raw == nil {
						nd := &osmNodeCount{nodeID, 1}
						wtTxn.Insert(edgeTable, nd)
					} else {

						newCnt := raw.(*osmNodeCount).Count + 1
						nd := &osmNodeCount{nodeID, newCnt}
						wtTxn.Insert(edgeTable, nd)
					}
				}

				wtTxn.Commit()
			}
		}
	}
}

func (e *edgeQuery) IsEdge (nd int64) bool {
	rdTxn := e.db.Txn(false)

	raw, _ := rdTxn.First(edgeTable, "id", nd)
	rdTxn.Abort()

	if raw != nil {
		if raw.(*osmNodeCount).Count > 1 {
			return true
		}
	}

	return false
}
