package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/go-memdb"
	"github.com/qedus/osmpbf"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/dijkstra"
)

type Node struct {
	OldID int64
	NewID int
}

type NodeUsage struct {
	ID int64
	Count int
}

// dijkstraCommand defines the version command.
func dijkstraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dijkstra",
		Short: "a dijkstra test",
		Long:  "a dijkstra test",
		Run: func(cmd *cobra.Command, args []string) {
			graph := dijkstra.NewGraph()

			schema := &memdb.DBSchema{
				Tables: map[string]*memdb.TableSchema{
					"node_usage": &memdb.TableSchema{
						Name: "node_usage",
						Indexes: map[string]*memdb.IndexSchema{
							"id": &memdb.IndexSchema{
								Name:    "id",
								Unique:  true,
								Indexer: &memdb.IntFieldIndex{Field: "ID"},
							},
							"count": &memdb.IndexSchema{
								Name:    "count",
								Unique:  true,
								Indexer: &memdb.IntFieldIndex{Field: "Count"},
							},
						},
					},
					"node": &memdb.TableSchema{
						Name: "node",
						Indexes: map[string]*memdb.IndexSchema{
							"id": &memdb.IndexSchema{
								Name:    "id",
								Unique:  true,
								Indexer: &memdb.IntFieldIndex{Field: "OldID"},
							},
							"new_id": &memdb.IndexSchema{
								Name:    "new_id",
								Unique:  true,
								Indexer: &memdb.IntFieldIndex{Field: "NewID"},
							},
						},
					},
				},
			}

			// Create a new data base
			db, err := memdb.NewMemDB(schema)
			if err != nil {
				panic(err)
			}

			//"/Users/thanoskontos/Downloads/greece_for_routing.osm.pbf"
			OSMFilename := "/Users/thanoskontos/Downloads/bremen_for_routing.osm.pbf"

			logNodes(db, OSMFilename)
			fmt.Println("Done logging nodes")


			f, err := os.Open(OSMFilename)
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

			var autoInc = 0
			var lastAddedVertex = 0
			var nc, wc, rc uint64
			for {
				if v, err := d.Decode(); err == io.EOF {
					break
				} else if err != nil {
					log.Fatal(err)
				} else {
					switch v := v.(type) {
					case *osmpbf.Node:
						nc++
					case *osmpbf.Way:
						var intersections []int64
						for i, nd := range v.NodeIDs {
							if i == 0 || i == len(v.NodeIDs)-1 {
								intersections = append(intersections, nd)
								continue
							}

							rdTxn := db.Txn(false)

							raw, _ := rdTxn.First("node_usage", "id", nd)
							rdTxn.Abort()

							if raw != nil {
								if raw.(*NodeUsage).Count > 1 {
									intersections = append(intersections, raw.(*NodeUsage).ID)
								}
							}
						}

						// Now reform the intersections to new ids
						wtTxn := db.Txn(true)

						var newIntersectionIDs []int
						for _, isnNd := range intersections {
							rdTxn := db.Txn(false)
							raw, _ := rdTxn.First("node", "id", isnNd)
							rdTxn.Abort()

							if raw == nil {
								autoInc++
								nd := &Node{isnNd, autoInc}
								wtTxn.Insert("node", nd)

								newIntersectionIDs = append(newIntersectionIDs, autoInc)
							} else {
								newIntersectionIDs = append(newIntersectionIDs, raw.(*Node).NewID)
							}
						}

						wtTxn.Commit()

						vtx := addIntersectionsToGraph(graph, newIntersectionIDs, lastAddedVertex)
						if vtx != -1 {
							lastAddedVertex = vtx
						}

						// Process Way v.

						wc++
					case *osmpbf.Relation:
						// Process Relation v.
						rc++
					default:
						log.Fatalf("unknown type %T\n", v)
					}
				}
			}

			fmt.Println("Nodes: %d, Ways: %d, Relations: %d\n", nc, wc, rc)

			//rdTxn := db.Txn(false)
			//raw, _ := rdTxn.First("node", "id", 26171771)
			//rdTxn.Abort()
			//if raw != nil {
			//	fmt.Println(raw.(*Node).NewID)
			//}
			//
			//rdTxn = db.Txn(false)
			//raw, _ = rdTxn.First("node", "id", 26207142)
			//rdTxn.Abort()
			//if raw != nil {
			//	fmt.Println(raw.(*Node).NewID)
			//}

			best, err := graph.Shortest(2173, 2201)

			fmt.Println("Shortest distance ", best.Distance, " following path ", best.Path)
		},
	}
}

func addIntersectionsToGraph(graph *dijkstra.Graph, intersections []int, previousLastAddedVertex int) int {
	previous := 0
	lastAddedVertex := -1

	for i, isn := range intersections {
		if isn > previousLastAddedVertex || previousLastAddedVertex == 0 {
			//fmt.Println("added vertex", isn, "with previously added vertex", previousLastAddedVertex)
			graph.AddVertex(isn)

			lastAddedVertex = isn
		}

		if i == 0 {
			previous = isn
		} else {
			//fmt.Println(fmt.Sprintf("graph.AddArc(%d, %d, 1)", isn, previous))

			graph.AddArc(isn, previous, 1)
			graph.AddArc(previous, isn, 1)
		}
	}

	//fmt.Println(lastAddedVertex)

	return lastAddedVertex
}

func logNodes(db *memdb.MemDB, filename string) {
	f, err := os.Open(filename)
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
				wtTxn := db.Txn(true)

				for _, nodeID := range v.NodeIDs {
					rdTxn := db.Txn(false)
					raw, err := rdTxn.First("node_usage", "id", nodeID)
					rdTxn.Abort()

					if err == nil && raw == nil {
						nd := &NodeUsage{nodeID, 1}
						wtTxn.Insert("node_usage", nd)
					} else {


						newCnt := raw.(*NodeUsage).Count + 1
						nd := &NodeUsage{nodeID, newCnt}
						wtTxn.Insert("node_usage", nd)
					}
				}

				wtTxn.Commit()
			}
		}
	}
}
