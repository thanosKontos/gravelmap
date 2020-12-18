package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thanosKontos/gravelmap/log"
)

func TestImportCreatesAllFiles(t *testing.T) {
	defer os.RemoveAll("../fixtures/_files")

	s := NewImport(ImportConfig{
		OsmFilemame:          "../fixtures/sample.osm.pbf",
		OutputDir:            "../fixtures/_files",
		ElevationDir:         "../fixtures",
		Log:                  log.NewNullLog(),
		Osm2GmUseFilesystem:  false,
		ProfileName:          "abc",
		ProfileFilename:      fmt.Sprintf("../profiles/mtb.yaml"),
		ElevationCredentials: ElevationCredentials{Username: "", Password: ""},
	})

	err := s.Import()
	assert.Nil(t, err)

	filenames := []string{"edge_start.bin", "edge_to_polylines_lookup.bin", "polylines.bin", "graph_abc.gob", "edge_bbox/N37E23_9.bin", "edge_bbox/N38E23_0.bin"}
	for _, filename := range filenames {
		assert.FileExists(t, fmt.Sprintf("../fixtures/_files/%s", filename))

		f, _ := os.Stat(fmt.Sprintf("../fixtures/_files/%s", filename))
		assert.NotEmpty(t, f.Size())
	}
}
