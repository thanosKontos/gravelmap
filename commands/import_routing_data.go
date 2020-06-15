package commands

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/elevation/hgt"
	"github.com/thanosKontos/gravelmap/encode"
	"github.com/thanosKontos/gravelmap/graph"
	"github.com/thanosKontos/gravelmap/node"
	"github.com/thanosKontos/gravelmap/node2point"
	"github.com/thanosKontos/gravelmap/osm"
	"github.com/thanosKontos/gravelmap/path"
	"github.com/thanosKontos/gravelmap/routing_algorithm/dijkstra"
	"github.com/thanosKontos/gravelmap/way"
)

// importRoutingDataCommand imports data from an OSM file.
func importRoutingDataCommand() *cobra.Command {
	var (
		inputFilename, routingMd string
		useFilesystem            bool
	)

	importRoutingDataCmd := &cobra.Command{
		Use:   "import-routing-data",
		Short: "import routing data",
		Long:  "import routing data",
	}

	importRoutingDataCmd.Flags().StringVar(&inputFilename, "input", "", "The osm input file.")
	importRoutingDataCmd.Flags().StringVar(&routingMd, "routing-mode", "bicycle", "The routing mode.")
	importRoutingDataCmd.Flags().BoolVar(&useFilesystem, "use-filesystem", false, "Use filesystem if your system runs out of memory (e.g. you are importing a large osm file on a low mem system). Will make import very slow!")

	importRoutingDataCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return importRoutingDataCmdRun(inputFilename, routingMd, useFilesystem)
	}

	return importRoutingDataCmd
}

type routingMode struct {
	graphFilename string
	weighter      gravelmap.Weighter
}

func importRoutingDataCmdRun(inputFilename string, routingMd string, useFilesystem bool) error {
	os.Mkdir("_files", 0777)

	// ## 1. Initially extract only the way nodes and keep them in a DB. Also keeps the GM identifier ##
	var osm2GmStore gravelmap.Osm2GmNodeReaderWriter
	if useFilesystem {
		osm2GmStore = node.NewOsm2GmNodeFileStore("_files")
	} else {
		osm2GmStore = node.NewOsm2GmNodeMemoryStore()
	}

	osm2GmNode := osm.NewOsmWayProcessor(inputFilename, osm2GmStore)
	err := osm2GmNode.Process()
	if err != nil {
		logger.Error(err)
		os.Exit(0)
	}
	logger.Info("Done preparing node in-memory DB")

	// ## 2. Store nodes to lookup files (nodeId -> lat/lon and lat/lon to closest nodeId)
	bboxFS := node2point.NewNodePointBboxFileStore("_files")
	osm2LatLngStore := node.NewOsm2LatLngMemoryStore()
	ndFileStore := node.NewOsmNodeProcessor(inputFilename, osm2GmStore, bboxFS, osm2LatLngStore)
	err = ndFileStore.Process()
	if err != nil {
		logger.Error(err)
		os.Exit(0)
	}
	logger.Info("Node file written")

	// ## 3. Process OSM ways (store way info and create graph)
	elevationGetterCloser := hgt.NewHgt("/tmp", os.Getenv("NASA_USERNAME"), os.Getenv("NASA_PASSWORD"), logger)
	distanceCalculator := distance.NewHaversine()
	pathSimplifier := path.NewSimplifiedDouglasPeucker(distanceCalculator)
	pointEncoder := encode.NewGooglemaps()

	var rms = map[string]routingMode{
		"bicycle": {"graph_bicycle.gob", way.NewBicycleWeight()},
		"foot":    {"graph_foot.gob", way.NewHikingWeight()},
	}
	routingMode := rms[routingMd]

	costEvaluator := way.NewCostEvaluate(distanceCalculator, elevationGetterCloser, routingMode.weighter)

	graph := graph.NewGraph()
	wayStorer := way.NewFileStore("_files", pointEncoder)
	wayAdderGetter := osm.NewOsm2GmWays(osm2GmStore, osm2LatLngStore, costEvaluator, pathSimplifier)

	osmWayFileRead := osm.NewOsmWayFileRead(inputFilename, wayStorer, graph, wayAdderGetter)
	err = osmWayFileRead.Process()
	if err != nil {
		logger.Error(err)
		os.Exit(0)
	}
	logger.Info("Ways processed")

	elevationGetterCloser.Close()

	// also persist it to file
	graphFile, _ := os.Create(fmt.Sprintf("_files/%s", routingMode.graphFilename))
	dataEncoder := gob.NewEncoder(graphFile)
	dataEncoder.Encode(graph)
	graphFile.Close()
	logger.Info("Graph created")

	dijkstra := dijkstra.NewDijkstra(graph)
	best, _ := dijkstra.FindShortest(14827, 1037)

	logger.Info(fmt.Sprintf("Shortest distance %d following path %#v", best.Distance, best.Path))

	return nil
}
