package prepare

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/go-memdb"
	"github.com/qedus/osmpbf"
)

const nodeTable = "gravelmap_osm_gm_node"

type nodeQuery struct {
	osmFilename  string
	db *memdb.MemDB
}

type Node struct {
	OldID int64
	NewID int
	Occurences int
}

func NewNodeQuerer(osmFilename string) *nodeQuery {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			nodeTable: {
				Name: nodeTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "OldID"},
					},
					"new_id": {
						Name:    "new_id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "NewId"},
					},
				},
			},
		},
	}

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}

	return &nodeQuery{
		osmFilename:  osmFilename,
		db: db,
	}
}

func (n *nodeQuery) Prepare () *memdb.MemDB {
	f, err := os.Open(n.osmFilename)
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

	inc := 0
	var wc uint64
	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		} else {
			switch v := v.(type) {
			case *osmpbf.Way:
				for _, osmNdID := range v.NodeIDs {
					ndDB := n.getOSMNodeIDFromDB(osmNdID)

					if ndDB == nil {
						inc++
						if osmNdID == 20974186 {
							fmt.Println(osmNdID, inc, 1)
						}

						wtTxn := n.db.Txn(true)
						wtTxn.Insert(nodeTable, &Node{osmNdID, inc, 1})
						wtTxn.Commit()
					} else {
						newCnt := ndDB.Occurences + 1

						if osmNdID == 20974186 {
							fmt.Println(ndDB.OldID, ndDB.NewID, newCnt)
						}

						wtTxn := n.db.Txn(true)
						wtTxn.Insert(nodeTable, &Node{ndDB.OldID, ndDB.NewID, newCnt})
						wtTxn.Commit()
					}
				}

				wc++
			}
		}
	}

	return n.db
}

func (n *nodeQuery) getOSMNodeIDFromDB(osmNdID int64) *Node {
	rdTxn := n.db.Txn(false)
	rs, _ := rdTxn.First(nodeTable, "id", osmNdID)
	rdTxn.Abort()

	if osmNdID == 20974186 {
		fmt.Println("Value read", rs)
	}

	if rs == nil {
		return nil
	} else {
		return rs.(*Node)
	}
}

