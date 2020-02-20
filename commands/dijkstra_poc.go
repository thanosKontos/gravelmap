package commands

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
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


			var latLngs []maps.LatLng
			for _, pathNd := range best.Path {
				test_node, _ := ndFileStore.Read(int32(pathNd))

				latLngs = append(latLngs, maps.LatLng{Lat: test_node.Lat, Lng: test_node.Lng})
			}
			fmt.Println(maps.Encode(latLngs))
		},
	}
}
