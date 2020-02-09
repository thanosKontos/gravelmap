package commands

import (
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/osmfilter"
	"log"
)

// createRoutingDBCommand defines the create route command.
func createFilterOSMCommand() *cobra.Command {
	var (
		inputFilename  string
		outputFilename string
	)

	createFilterOSMCmd := &cobra.Command{
		Use:   "filter-osm",
		Short: "filter osm",
		Long:  "filter osm to include just useful routing data",
	}

	createFilterOSMCmd.Flags().StringVar(&inputFilename, "input", "", "The osm input file.")
	createFilterOSMCmd.Flags().StringVar(&outputFilename, "output", "", "The osm input file.")
	createFilterOSMCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createfilterOSMCmdRun(inputFilename, outputFilename)
	}

	return createFilterOSMCmd
}

// createRoutingDBCmdRun defines the command run actions.
func createfilterOSMCmdRun(inputFilename, outputFilename string) error {
	if inputFilename == "" || outputFilename == "" {
		log.Fatalln("please specify input and output files")
	}

	osmium := osmfitler.NewOsmium(inputFilename, outputFilename)
	err := osmium.Filter()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("OSM file prepared.")

	return nil
}
