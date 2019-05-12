package commands

import (
	"github.com/spf13/cobra"
	osmium "github.com/thanosKontos/gravelmap/osmfilter"
	routing_import "github.com/thanosKontos/gravelmap/routing/osm2pgrouting/import"
	routing_prepare "github.com/thanosKontos/gravelmap/routing/pgrouting/prepare"
	"log"
	"os"
)

// createRoutingDataCommand defines the create routing command.
func createRoutingDataCommand() *cobra.Command {
	var (
		inputFilename string
		tagCostConf   string
	)

	createRoutingDataCmd := &cobra.Command{
		Use:   "add-data",
		Short: "add data to routing database",
		Long:  "manipulate osm data, extract tags and insert to routing database",
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
	osmiumFilter := osmium.NewOsmium(inputFilename, "/tmp/filtered.osm")
	err := osmiumFilter.Filter()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("OSM data filtered successfully.")

	pgPreparator, err := routing_prepare.NewPgRoutingPrep(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"), os.Getenv("DBDEFAULTDBNAME"))
	if err != nil {
		log.Fatal(err)
	}
	defer pgPreparator.Close()

	err = pgPreparator.Prepare()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database prepared.")

	pgImporter := routing_import.NewOsm2PgRouting(os.Getenv("DBUSER"), os.Getenv("DBPASS"), os.Getenv("DBNAME"), os.Getenv("DBPORT"), "/tmp/filtered.osm", tagCostConf)
	err = pgImporter.Import()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database filled.")

	return nil
}
