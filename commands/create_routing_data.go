package commands

import (
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/cli"
	"github.com/thanosKontos/gravelmap/osmfilter"
	"github.com/thanosKontos/gravelmap/routing_import"
	"github.com/thanosKontos/gravelmap/routing_prepare"
	"log"
	"os"
)

// createRoutingDataCommand defines the create route command.
func createRoutingDataCommand() *cobra.Command {
	var (
		inputFilename string
		tagCostConf   string
	)

	createRoutingDataCmd := &cobra.Command{
		Use:   "add-data",
		Short: "add data to route database",
		Long:  "manipulate osm data, extract tags and insert to route database",
	}

	createRoutingDataCmd.Flags().StringVar(&inputFilename, "input", "", "The osm input file.")
	createRoutingDataCmd.Flags().StringVar(&tagCostConf, "tag-cost-config", "", "Tag cost config.")

	createRoutingDataCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createRoutingDataCmdRun(inputFilename, tagCostConf)
	}

	return createRoutingDataCmd
}

// createRoutingDataCmdRun defines the command run actions.
func createRoutingDataCmdRun(inputFilename, tagCostConf string) error {
	// Filter useless osm tags
	osmiumFilter := osmfitler.NewOsmium(inputFilename, "/tmp/filtered.osm")
	err := osmiumFilter.Filter()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OSM data filtered successfully.")

	gmPreparer, err := routing_prepare.NewGravelmapPreparer(
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPORT"),
		os.Getenv("DBDEFAULTDBNAME"),
		cli.NewCLI(),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer gmPreparer.Close()

	err = gmPreparer.Prepare()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database prepared.")

	pgImporter := routing_import.NewOsm2PgRouting(
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPORT"),
		"/tmp/filtered.osm",
		tagCostConf,
	)
	err = pgImporter.Import()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database filled.")

	return nil
}
