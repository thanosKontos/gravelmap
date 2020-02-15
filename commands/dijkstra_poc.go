package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"encoding/gob"

	"github.com/hashicorp/go-memdb"
	"github.com/qedus/osmpbf"
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/dijkstra"
)

type EdgeNode struct {
	OsmNdID int64
	GmNdID  int
}

type OsmNodeCount struct {
	ID    int64
	Count int
}

// dijkstraCommand defines the version command.
func dijkstraPocCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "dijkstra-poc",
		Short: "a dijkstra test",
		Long:  "a dijkstra test",
		Run: func(cmd *cobra.Command, args []string) {

			//var latlngs = []maps.LatLng{{Lat: 39.87709, Lng: 32.74713}, {Lat: 39.87709, Lng: 32.74787}, {Lat: 39.87653, Lng: 32.74746}}
			//encoded := maps.Encode(latlngs)

			//decode, _ := maps.DecodePolyline("ynkrFq|zfE?sCnBrA")
			//fmt.Println(decode)
			//os.Exit(0)



			graph := dijkstra.NewGraph()

			schema := &memdb.DBSchema{
				Tables: map[string]*memdb.TableSchema{
					"osm_node_count": &memdb.TableSchema{
						Name: "osm_node_count",
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
					"edge_node": &memdb.TableSchema{
						Name: "edge_node",
						Indexes: map[string]*memdb.IndexSchema{
							"id": &memdb.IndexSchema{
								Name:    "id",
								Unique:  true,
								Indexer: &memdb.IntFieldIndex{Field: "OsmNdID"},
							},
							"new_id": &memdb.IndexSchema{
								Name:    "new_id",
								Unique:  true,
								Indexer: &memdb.IntFieldIndex{Field: "GmNdID"},
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

			//OSMFilename := "/Users/thanoskontos/Downloads/greece_for_routing.osm.pbf"
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

							raw, _ := rdTxn.First("osm_node_count", "id", nd)
							rdTxn.Abort()

							if raw != nil {
								if raw.(*OsmNodeCount).Count > 1 {
									intersections = append(intersections, raw.(*OsmNodeCount).ID)
								}
							}
						}

						// Now reform the intersections to new ids
						wtTxn := db.Txn(true)

						var newIntersectionIDs []int
						for _, isnNd := range intersections {
							rdTxn := db.Txn(false)
							raw, _ := rdTxn.First("edge_node", "id", isnNd)
							rdTxn.Abort()

							if raw == nil {
								autoInc++
								nd := &EdgeNode{isnNd, autoInc}
								wtTxn.Insert("edge_node", nd)

								newIntersectionIDs = append(newIntersectionIDs, autoInc)
							} else {
								newIntersectionIDs = append(newIntersectionIDs, raw.(*EdgeNode).GmNdID)
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

			dataFile, err := os.Create("_files/graph.gob")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			dataEncoder := gob.NewEncoder(dataFile)
			dataEncoder.Encode(&graph)

			dataFile.Close()

			//rdTxn := db.Txn(false)
			//raw, _ := rdTxn.First("edge_node", "id", 26171771)
			//rdTxn.Abort()
			//if raw != nil {
			//	fmt.Println(raw.(*EdgeNode).GmNdID)
			//}
			//
			//rdTxn = db.Txn(false)
			//raw, _ = rdTxn.First("edge_node", "id", 26207142)
			//rdTxn.Abort()
			//if raw != nil {
			//	fmt.Println(raw.(*EdgeNode).GmNdID)
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
					raw, err := rdTxn.First("osm_node_count", "id", nodeID)
					rdTxn.Abort()

					if err == nil && raw == nil {
						nd := &OsmNodeCount{nodeID, 1}
						wtTxn.Insert("osm_node_count", nd)
					} else {

						newCnt := raw.(*OsmNodeCount).Count + 1
						nd := &OsmNodeCount{nodeID, newCnt}
						wtTxn.Insert("osm_node_count", nd)
					}
				}

				wtTxn.Commit()
			}
		}
	}
}
