package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/service"
)

// importRoutingDataCommand imports data from an OSM file.
func importRoutingDataCommand() *cobra.Command {
	var (
		inputFilename, profileName string
		useFilesystem              bool
	)

	importRoutingDataCmd := &cobra.Command{
		Use:   "import-routing-data",
		Short: "import routing data",
		Long:  "import routing data",
	}

	importRoutingDataCmd.Flags().StringVar(&inputFilename, "input", "", "The osm input file.")
	importRoutingDataCmd.Flags().StringVar(&profileName, "profile", "mtb", "The profile (routing type) to use.")
	importRoutingDataCmd.Flags().BoolVar(&useFilesystem, "use-filesystem", false, "Use filesystem if your system runs out of memory (e.g. you are importing a large osm file on a low mem system). Will make import very slow!")

	importRoutingDataCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return importRoutingDataCmdRun(inputFilename, profileName, useFilesystem)
	}

	return importRoutingDataCmd
}

func importRoutingDataCmdRun(inputFilename string, profileName string, useFilesystem bool) error {
	s := service.NewImport(service.ImportConfig{
		OsmFilemame:          inputFilename,
		OutputDir:            "_files",
		ElevationDir:         "/tmp",
		Log:                  logger,
		Osm2GmUseFilesystem:  useFilesystem,
		ProfileName:          profileName,
		ProfileFilename:      fmt.Sprintf("profiles/%s.yaml", profileName),
		ElevationCredentials: service.ElevationCredentials{Username: os.Getenv("NASA_USERNAME"), Password: os.Getenv("NASA_PASSWORD")},
	})

	err := s.Import()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}

	return nil
}
