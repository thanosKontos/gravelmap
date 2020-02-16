package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/node_db"
	"github.com/thanosKontos/gravelmap/prepare"
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

			//var latlngs = []maps.LatLng{{Lat: 39.87709, Lng: 32.74713}, {Lat: 39.87709, Lng: 32.74787}, {Lat: 39.87653, Lng: 32.74746}}
			//encoded := maps.Encode(latlngs)

			//decode, _ := maps.DecodePolyline("ynkrFq|zfE?sCnBrA")
			//fmt.Println(decode)
			//os.Exit(0)

			ndDB := node_db.NewNodeMapDB()

			nodeQuery:= prepare.NewNodeQuerer(OSMFilename, ndDB)
			nodeDB := nodeQuery.Prepare()

			fmt.Println("done 1")

			gmGraph := prepare.NewGraph(OSMFilename, nodeDB)
			gmGraph.Prepare()

			fmt.Println("done 2")
			os.Exit(0)
		},
	}
}
