package commands

import (
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/elevation/srtm_ascii"
	"log"
	"os"
)

// createRoutingDataCommand defines the create route command.
func importElevationCommand() *cobra.Command {
	var (
		inputFilename string
	)

	createRoutingDataCmd := &cobra.Command{
		Use:   "import-elevation",
		Short: "import elevation data to route DB",
		Long:  "import elevation data to route DB",
	}

	createRoutingDataCmd.Flags().StringVar(&inputFilename, "input", "", "The asc input file.")

	createRoutingDataCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return importElevationCmdRun(inputFilename)
	}

	return createRoutingDataCmd
}

// createRoutingDataCmdRun defines the command run actions.
func importElevationCmdRun(inputFilename string) error {
	srtm, err := srtm_ascii.NewSRTM(
		inputFilename,
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPORT"),
		logger,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = srtm.Import()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Elevation data imported")

	return nil
}
