package commands

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/edge"
	"github.com/thanosKontos/gravelmap/elevation"
	"github.com/thanosKontos/gravelmap/node"
	"github.com/thanosKontos/gravelmap/prepare"
	"github.com/thanosKontos/gravelmap/way"
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
			os.Mkdir("_files", 0777)

			//OSMFilename := "/Users/thanoskontos/Downloads/greece_for_routing.osm.pbf"
			OSMFilename := "/Users/thanoskontos/Downloads/bremen_for_routing.osm.pbf"

			// ## 1. Initially extract only the way nodes and keep them in a DB. Also keeps the GM identifier ##
			osm2GmStore := node.NewOsm2GmNodeMemoryStore()
			osm2GmNode:= prepare.NewOsm2GmNode(OSMFilename, osm2GmStore)
			osm2GmNode.Extract()

			logger.Info("Done preparing node in-memory DB")

			// ## 2. Store nodes to lookup files (nodeId -> lat/lon and lat/lon to closest nodeId)
			bboxFS := edge.NewBBoxFileStore("_files")
			ndFileStore := node.NewNodeFileStore("_files", OSMFilename, osm2GmStore, bboxFS)
			err := ndFileStore.Persist()
			if err != nil {
				logger.Error(err)
				os.Exit(0)
			}
			logger.Info("Node file written")


			// TODO: here we need to pass to the graph preperator some kind of elevation grader and distance calculator

			// ## 3. Create the dijkstra graph that we will use to do the actual routing ##
			elevationGetterCloser := elevation.NewHgt("/tmp", logger)
			distanceCalculator := distance.NewHaversine()
			costEvaluator := way.NewCostEvaluate(distanceCalculator, elevationGetterCloser)
			gmGraph := prepare.NewGraph(OSMFilename, osm2GmStore, costEvaluator)
			gmGraph.Prepare()

			elevationGetterCloser.Close()

			// also persist it to file
			graphFile, _ := os.Create("_files/graph.gob")
			dataEncoder := gob.NewEncoder(graphFile)
			dataEncoder.Encode(gmGraph.GetGraph())
			graphFile.Close()
			logger.Info("Graph created")

			// ## 4. Store polylines for ways
			wayFileStore := way.NewWayFileStore("_files", OSMFilename, osm2GmStore, ndFileStore)
			err = wayFileStore.Persist()
			if err != nil {
				logger.Error("Way files written")
				os.Exit(0)
			}
			logger.Info("Way files written")



			dGraph := gmGraph.GetGraph()
			//best, _ := dGraph.Shortest(1, 2173)
			best, _ := dGraph.Shortest(14827, 1037)

			logger.Info(fmt.Sprintf("Shortest distance %d following path %#v", best.Distance, best.Path))
		},
	}
}
