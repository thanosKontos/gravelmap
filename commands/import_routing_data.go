package commands

import (
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/routing_import"
	"log"
	"os"
)

// importRoutingDataCommand defines the create route command.
func importRoutingDataCommand() *cobra.Command {
	var (
		inputFilename string
		tagCostConf   string
	)

	importRoutingDataCmd := &cobra.Command{
		Use:   "import-osm",
		Short: "import data to routing database",
		Long:  "manipulate osm data, extract tags and insert to route database",
	}

	importRoutingDataCmd.Flags().StringVar(&inputFilename, "input", "", "The osm input file.")
	importRoutingDataCmd.Flags().StringVar(&tagCostConf, "tag-cost-config", "", "Tag cost config.")

	importRoutingDataCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return importRoutingDataCmdRun(inputFilename, tagCostConf)
	}

	return importRoutingDataCmd
}

// importRoutingDataCmdRun defines the command run actions.
func importRoutingDataCmdRun(inputFilename, tagCostConf string) error {
	pgImporter := routing_import.NewOsm2PgRouting(
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPORT"),
		inputFilename,
		tagCostConf,
	)
	err := pgImporter.Import()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("OSM imported to DB.")

	return nil
}
