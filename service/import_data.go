package service

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

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

type importService struct {
	conf ImportConfig
}

type ElevationCredentials struct {
	Username string
	Password string
}

type ImportConfig struct {
	OsmFilemame          string
	OutputDir            string
	ElevationDir         string
	Log                  gravelmap.Logger
	Osm2GmUseFilesystem  bool
	ProfileName          string
	ProfileFilename      string
	ElevationCredentials ElevationCredentials
}

func NewImport(conf ImportConfig) importService {
	return importService{conf: conf}
}

func (i importService) Import() error {
	os.Mkdir(i.conf.OutputDir, 0777)

	// ## 1. Initially extract only the way nodes and keep them in a DB. Also keeps the GM identifier ##
	var osm2GmStore gravelmap.Osm2GmNodeReaderWriter
	if i.conf.Osm2GmUseFilesystem {
		osm2GmStore = node.NewOsm2GmNodeFileStore(i.conf.OutputDir)
	} else {
		osm2GmStore = node.NewOsm2GmNodeMemoryStore()
	}

	osm2GmNode := osm.NewOsmWayProcessor(i.conf.OsmFilemame, osm2GmStore)
	err := osm2GmNode.Process()
	if err != nil {
		return err
	}
	i.conf.Log.Info("Done preparing node in-memory DB")

	// ## 2. Store nodes to lookup files (nodeId -> lat/lon and lat/lon to closest nodeId)
	bboxFS := node2point.NewNodePointBboxFileStore(i.conf.OutputDir)
	osm2LatLngStore := node.NewOsm2LatLngMemoryStore()
	ndFileStore := osm.NewOsmNodeProcessor(i.conf.OsmFilemame, osm2GmStore, bboxFS, osm2LatLngStore)
	err = ndFileStore.Process()
	if err != nil {
		return err
	}
	i.conf.Log.Info("Node file written")

	// ## 3. Process OSM ways (store way info and create graph)
	elevationGetterCloser := hgt.NewNasaHgt(
		i.conf.ElevationDir,
		i.conf.ElevationCredentials.Username,
		i.conf.ElevationCredentials.Password,
		i.conf.Log,
	)
	distanceCalculator := distance.NewHaversine()

	weightConf := way.WeightConfig{}
	yamlFile, err := ioutil.ReadFile(i.conf.ProfileFilename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, &weightConf)
	if err != nil {
		return err
	}

	pathEncoder := path.NewGooglePolyline()
	wayStorer := way.NewFileStore(i.conf.OutputDir, pathEncoder)
	pathSimplifier := path.NewSimpleSimplifiedPath(distanceCalculator)
	costEvaluator := way.NewCostEvaluate(distanceCalculator, elevationGetterCloser, way.NewDefaultWeight(weightConf))
	wayAdderGetter := osm.NewOsm2GmWays(osm2GmStore, osm2LatLngStore, costEvaluator, pathSimplifier)

	graph := graph.NewWeightedBidirectionalGraph()
	osmWayFileRead := osm.NewOsmWayFileRead(i.conf.OsmFilemame, wayStorer, graph, wayAdderGetter)
	err = osmWayFileRead.Process()
	if err != nil {
		return err
	}
	i.conf.Log.Info("Ways processed")

	elevationGetterCloser.Close()

	// also persist it to file
	graphFile, err := os.Create(fmt.Sprintf("%s/graph_%s.gob", i.conf.OutputDir, i.conf.ProfileName))
	if err != nil {
		return err
	}
	dataEncoder := gob.NewEncoder(graphFile)
	dataEncoder.Encode(graph)
	graphFile.Close()
	i.conf.Log.Info("Graph created")

	return nil
}
