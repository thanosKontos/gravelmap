package commands

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/node"
	"github.com/thanosKontos/gravelmap/prepare"
	"github.com/thanosKontos/gravelmap/way"
	"googlemaps.github.io/maps"
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

			testWays := []int32{1,86123,135138,135133,121181,85173,5519,121174,116378,85694,86138,63143,4689,85131,121195,86120,85760,112247,63577,112242,112237,135110,56424,85141,135102,56428,135077,973,132006,135067,82937,698,85158,135054,132060,13339,132055,107433,112180,134875,85124,115145,39299,132050,96635,138762,138765,152794,112184,152793,152792,1690,86599,86594,86592,86591,86581,121805,2170,2173}

			var testWayPairs []gravelmap.Way
			var prev int32 = 0
			for i, testway := range testWays {
				if i == 0 {
					prev = testway
					continue
				}

				testWayPairs = append(testWayPairs, gravelmap.Way{prev, testway})

				prev = testway
			}

			wayFile := way.NewWayFileRead("_files")
			polylines, _ := wayFile.Read(testWayPairs)

			//fmt.Println(polylines)

			var latLngs []maps.LatLng
			for _, pl := range polylines {
				tmpLatLngs, _ := maps.DecodePolyline(pl)
				for _, latlng := range tmpLatLngs {
					latLngs = append(latLngs, maps.LatLng{Lat: latlng.Lat, Lng: latlng.Lng})
				}
			}

			fmt.Println(maps.Encode(latLngs))
			os.Exit(0)



			//OSMFilename := "/Users/thanoskontos/Downloads/greece_for_routing.osm.pbf"
			OSMFilename := "/Users/thanoskontos/Downloads/bremen_for_routing.osm.pbf"

			osm2GmStore := node.NewOsm2GmNodeMemoryStore()
			nodeQuery:= prepare.NewOsm2GmNodeExtractor(OSMFilename, osm2GmStore)
			nodeDB := nodeQuery.Extract()

			log.Println("Done preparing node in-memory DB")

			gmGraph := prepare.NewGraph(OSMFilename, nodeDB)
			gmGraph.Prepare()

			log.Println("Done creating graph")

			dGraph := gmGraph.GetGraph()
			best, _ := dGraph.Shortest(1, 2173)

			log.Println("Shortest distance", best.Distance, "following path", best.Path)

			ndFileStore := node.NewNodeFileStore("_files", OSMFilename, nodeDB)
			err := ndFileStore.Persist()
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Node file written")

			wayFileStore := way.NewWayFileStore("_files", OSMFilename, nodeDB, ndFileStore)
			err = wayFileStore.Persist()
			if err != nil {
				log.Fatal(err)
			}

			log.Println("Way files written")


			fmt.Println("------")
			//var latLngs []maps.LatLng
			for _, pathNd := range best.Path {
				fmt.Print(pathNd, ",")
				//test_node, _ := ndFileStore.Read(int32(pathNd))

				//latLngs = append(latLngs, maps.LatLng{Lat: test_node.Lat, Lng: test_node.Lng})
			}
			//fmt.Println(maps.Encode(latLngs))
		},
	}
}
