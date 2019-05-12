package osmium

import (
	"errors"
	"fmt"
	"os/exec"
)

type Osmium struct {
	inputFilename  string
	outputFilename string
}

func NewOsmium(inputFilename, outputFilename string) *Osmium {
	return &Osmium{
		inputFilename:  inputFilename,
		outputFilename: outputFilename,
	}
}

func (o *Osmium) Filter() error {
	cmd := exec.Command("osmium", "tags-filter", o.inputFilename, "w/highway", "-o", "/tmp/filtered_tmp.osm", "--overwrite")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(fmt.Sprintf("%s : %s", err, out))
	}

	cmd = exec.Command("osmium", "tags-filter", "-i", "/tmp/filtered_tmp.osm", "w/highway=motorway,trunk,motorway_link,trunk_link", "w/access=private", "-o", o.outputFilename, "--overwrite")
	out, err = cmd.CombinedOutput()

	if err != nil {
		return errors.New(fmt.Sprintf("%s : %s", err, out))
	}

	return nil
}