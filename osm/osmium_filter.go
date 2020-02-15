package osm

import (
	"errors"
	"fmt"
	"os/exec"
)

type osmium struct {
	inputFilename  string
	outputFilename string
}

func NewOsmium(inputFilename, outputFilename string) *osmium {
	return &osmium{
		inputFilename:  inputFilename,
		outputFilename: outputFilename,
	}
}

func (o *osmium) Filter() error {
	cmd := exec.Command("osmium", "tags-filter", o.inputFilename, "w/highway", "-o", "/tmp/filtered_tmp.osm.pbf", "--overwrite")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%s : %s", err, out))
	}

	cmd = exec.Command("osmium", "tags-filter", "-i", "/tmp/filtered_tmp.osm.pbf", "w/highway=motorway,trunk,motorway_link,trunk_link", "w/name:en=*", "w/access=private", "-o", o.outputFilename, "--overwrite")
	out, err = cmd.CombinedOutput()

	//sed -E '/name:el|name:en|name:de|name:ru|name:bg|name:es|name:fr|name:tr|name:sq|traffic_signals|traffic_sign|wikidata|wikipedia|traffic_calming|operator|grave_yard|public_transport|crossing_ref|crossing|created_by/d' /Users/thanoskontos/Downloads/greece_for_routing.osm > /Users/thanoskontos/Downloads/greece_for_routing2.osm

	if err != nil {
		return errors.New(fmt.Sprintf("%s : %s", err, out))
	}

	return nil
}
