package prepare

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/go-memdb"
	"github.com/qedus/osmpbf"
)

type edge struct {
	osmFilename  string
	db *memdb.MemDB
}

type osmNodeCount struct {
	ID    int64
	Count int
}

func NewEdge(osmFilename string) *edge {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"gravelmap_osm_node_count": {
				Name: "gravelmap_osm_node_count",
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

	return &edge{
		osmFilename:  osmFilename,
		db: db,
	}
}

func (e *edge) Prepare () {
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
					raw, err := rdTxn.First("gravelmap_osm_node_count", "id", nodeID)
					rdTxn.Abort()

					if err == nil && raw == nil {
						nd := &osmNodeCount{nodeID, 1}
						wtTxn.Insert("gravelmap_osm_node_count", nd)
					} else {

						newCnt := raw.(*osmNodeCount).Count + 1
						nd := &osmNodeCount{nodeID, newCnt}
						wtTxn.Insert("gravelmap_osm_node_count", nd)
					}
				}

				wtTxn.Commit()
			}
		}
	}
}

func (e *edge) IsEdge (nd int64) bool {
	rdTxn := e.db.Txn(false)

	raw, _ := rdTxn.First("gravelmap_osm_node_count", "id", nd)
	rdTxn.Abort()

	if raw != nil {
		if raw.(*osmNodeCount).Count > 1 {
			return true
		}
	}

	return false
}
