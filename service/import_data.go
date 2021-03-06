package service

import (
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
	if err := osm2GmNode.Process(); err != nil {
		return err
	}
	i.conf.Log.Info("Done preparing node in-memory DB")

	// ## 2. Store nodes to lookup files (nodeId -> lat/lon and lat/lon to closest nodeId)
	bboxFS := node2point.NewNodePointBboxFileStore(i.conf.OutputDir)
	osm2LatLngStore := node.NewOsm2LatLngMemoryStore()
	ndFileStore := osm.NewOsmNodeProcessor(i.conf.OsmFilemame, osm2GmStore, bboxFS, osm2LatLngStore)
	if err := ndFileStore.Process(); err != nil {
		return err
	}
	i.conf.Log.Info("Node file written")

	// ## 3. Process OSM ways (store way info and create graph)
	elevationFileStorage := hgt.NewMemcacheNasaElevationFileStorage(
		i.conf.ElevationDir,
		i.conf.ElevationCredentials.Username,
		i.conf.ElevationCredentials.Password,
		i.conf.Log,
	)
	elevationWayGetterCloser := hgt.NewHgt(elevationFileStorage, i.conf.Log)
	distanceCalculator := distance.NewHaversine()

	weightConf := way.WeightConfig{}
	yamlFile, err := ioutil.ReadFile(i.conf.ProfileFilename)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, &weightConf); err != nil {
		return err
	}

	pathEncoder := path.NewGooglePolyline()
	wayStorer := way.NewFileStore(i.conf.OutputDir, pathEncoder)
	pathSimplifier := path.NewSimpleSimplifiedPath(distanceCalculator)
	costEvaluator := way.NewCostEvaluate(distanceCalculator, elevationWayGetterCloser, way.NewDefaultWeight(weightConf))
	wayAdderGetter := osm.NewOsm2GmWays(osm2GmStore, osm2LatLngStore, costEvaluator, pathSimplifier)

	g := graph.NewWeightedBidirectionalGraph()
	osmWayFileRead := osm.NewOsmWayFileRead(i.conf.OsmFilemame, wayStorer, g, wayAdderGetter)
	if err = osmWayFileRead.Process(); err != nil {
		return err
	}
	i.conf.Log.Info("Ways processed")

	elevationFileStorage.Close()

	repo := graph.NewGobRepo(i.conf.OutputDir)
	if err = repo.Store(g, i.conf.ProfileName); err != nil {
		return err
	}
	i.conf.Log.Info("Graph created")

	return nil
}
