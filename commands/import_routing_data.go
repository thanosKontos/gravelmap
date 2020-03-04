package commands

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/edge"
	"github.com/thanosKontos/gravelmap/elevation"
	"github.com/thanosKontos/gravelmap/encode"
	graph2 "github.com/thanosKontos/gravelmap/graph"
	"github.com/thanosKontos/gravelmap/node"
	"github.com/thanosKontos/gravelmap/osm"
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

// importRoutingDataCommand imports data from an OSM file.
func importRoutingDataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "import-routing-data",
		Short: "import routing data",
		Long:  "a dijkstra test",
		Run: func(cmd *cobra.Command, args []string) {
			os.Mkdir("_files", 0777)

			//OSMFilename := "/Users/thanoskontos/Downloads/greece_for_routing.osm.pbf"
			OSMFilename := "/Users/thanoskontos/Downloads/bremen_for_routing.osm.pbf"
			//OSMFilename := "/Users/thanoskontos/Downloads/evzonas_for_routing.osm.pbf"

			// ## 1. Initially extract only the way nodes and keep them in a DB. Also keeps the GM identifier ##
			osm2GmStore := node.NewOsm2GmNodeMemoryStore()
			osm2GmNode := osm.NewOsmNodeFileRead(OSMFilename, osm2GmStore)
			osm2GmNode.Process()

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

			// ## 3. Process OSM ways (store way info and create graph)
			elevationGetterCloser := elevation.NewHgt("/tmp", os.Getenv("NASA_USERNAME"), os.Getenv("NASA_PASSWORD"), logger)
			distanceCalculator := distance.NewHaversine()
			weighter := way.NewBicycleWeight()

			costEvaluator := way.NewCostEvaluate(distanceCalculator, elevationGetterCloser, weighter)
			pointEncoder := encode.NewGooglemaps()

			graph := graph2.NewDijkstra()
			wayStorer := way.NewFileStore("_files", pointEncoder)
			wayAdderGetter := osm.NewOsm2GmWays(osm2GmStore, ndFileStore, costEvaluator)

			osmWayFileRead := osm.NewOsmWayFileRead(OSMFilename, wayStorer, graph, wayAdderGetter)
			err = osmWayFileRead.Process()
			if err != nil {
				logger.Error(err)
				os.Exit(0)
			}
			logger.Info("Ways processed")

			elevationGetterCloser.Close()

			// also persist it to file
			graphFile, _ := os.Create("_files/graph.gob")
			dataEncoder := gob.NewEncoder(graphFile)
			dataEncoder.Encode(graph.Get())
			graphFile.Close()
			logger.Info("Graph created")

			dGraph := graph.Get()
			best, _ := dGraph.Shortest(14827, 1037)

			logger.Info(fmt.Sprintf("Shortest distance %d following path %#v", best.Distance, best.Path))
		},
	}
}
