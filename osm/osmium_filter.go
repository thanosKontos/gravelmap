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

	cmd = exec.Command("osmium", "tags-filter", "-i", "/tmp/filtered_tmp.osm.pbf", "w/highway=motorway,trunk,motorway_link,trunk_link,raceway", "w/service=parking_aisle,drive-through,driveway", "w/access=private,customers", "-o", o.outputFilename, "--overwrite")
	out, err = cmd.CombinedOutput()

	if err != nil {
		return errors.New(fmt.Sprintf("%s : %s", err, out))
	}

	return nil
}
