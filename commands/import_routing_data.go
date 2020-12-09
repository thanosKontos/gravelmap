package commands

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/distance"
	"github.com/thanosKontos/gravelmap/elevation/hgt"
	"github.com/thanosKontos/gravelmap/graph"
	"github.com/thanosKontos/gravelmap/node"
	"github.com/thanosKontos/gravelmap/node2point"
	"github.com/thanosKontos/gravelmap/osm"
	"github.com/thanosKontos/gravelmap/path"
	"github.com/thanosKontos/gravelmap/way"
	"gopkg.in/yaml.v2"
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

var data = `
weight_offroad: 0.6

weight_vehicle_acceptance:
  exclusively: 0.7
  yes: 1.0
  partially: 2.0
  maybe: 10000000.0
  no: 10000000.0

vehicle_acceptance_tags:
  exclusively:
    tags:
      - mtb:scale
    values:
      - bicycle:
        - yes
        - permissive
        - designated
  no:
    values:
      - bicycle:
        - no
  maybe:
    values:
      - highway:
        - footway
        - path
`

type T struct {
	WeightOffroad             float64 `yaml:"weight_offroad"`
	Weight_vehicle_acceptance struct {
		Exclusively float64
		Yes         float64
		Partially   float64
		Maybe       float64
		No          float64
	}

	Vehicle_acceptance_tags struct {
		Exclusively struct {
			Tags   []string
			Values interface{}
		}
		Yes struct {
			Tags   []string
			Values interface{}
		}
		Partially struct {
			Tags   []string
			Values interface{}
		}
		Maybe struct {
			Tags   []string
			Values interface{}
		}
		No struct {
			Tags   []string
			Values interface{}
		}
	}
}

func importRoutingDataCmdRun(inputFilename string, routingMd string, useFilesystem bool) error {

	t := T{}

	errzzz := yaml.Unmarshal([]byte(data), &t)
	if errzzz != nil {
		log.Fatalf("error: %v", errzzz)
	}
	fmt.Println(t.WeightOffroad)
	fmt.Println(t.Weight_vehicle_acceptance)
	fmt.Println(t.Vehicle_acceptance_tags.Exclusively.Tags)
	fmt.Println(t.Vehicle_acceptance_tags.Exclusively.Values)
	fmt.Println(t.Vehicle_acceptance_tags.Yes.Tags)
	fmt.Println(t.Vehicle_acceptance_tags.Yes.Values)

	os.Exit(0)

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
	ndFileStore := osm.NewOsmNodeProcessor(inputFilename, osm2GmStore, bboxFS, osm2LatLngStore)
	err = ndFileStore.Process()
	if err != nil {
		logger.Error(err)
		os.Exit(0)
	}
	logger.Info("Node file written")

	// ## 3. Process OSM ways (store way info and create graph)
	elevationGetterCloser := hgt.NewNasaHgt("/tmp", os.Getenv("NASA_USERNAME"), os.Getenv("NASA_PASSWORD"), logger)
	distanceCalculator := distance.NewHaversine()

	var rms = map[string]routingMode{
		"bicycle": {"graph_bicycle.gob", way.NewBicycleWeight()},
		"foot":    {"graph_foot.gob", way.NewHikingWeight()},
	}
	routingMode := rms[routingMd]

	pathEncoder := path.NewGooglePolyline()
	wayStorer := way.NewFileStore("_files", pathEncoder)
	pathSimplifier := path.NewSimpleSimplifiedPath(distanceCalculator)
	costEvaluator := way.NewCostEvaluate(distanceCalculator, elevationGetterCloser, routingMode.weighter)
	wayAdderGetter := osm.NewOsm2GmWays(osm2GmStore, osm2LatLngStore, costEvaluator, pathSimplifier)

	graph := graph.NewWeightedBidirectionalGraph()
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

	return nil
}
