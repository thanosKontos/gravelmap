package commands

import (
	"database/sql"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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
	cmd := exec.Command("osmium", "tags-filter", inputFilename, "w/highway", "-o", "/tmp/filtered_tmp.osm", "--overwrite")
	_, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		return err
	}

	cmd = exec.Command("osmium", "tags-filter", "-i", "/tmp/filtered_tmp.osm", "w/highway=motorway,trunk,motorway_link,trunk_link", "w/access=private", "-o", "/tmp/filtered.osm", "--overwrite")
	_, err = cmd.CombinedOutput()

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("OSM data filtered successfully.")

	connStr := fmt.Sprintf(
		"user=%s dbname=%s password=%s port=%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBDEFAULTDBNAME"),
		os.Getenv("DBPASS"),
		os.Getenv("DBPORT"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("DROP DATABASE IF EXISTS " + os.Getenv("DBNAME"))
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("CREATE DATABASE " + os.Getenv("DBNAME"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Database created.")
	db.Close()

	connStr = fmt.Sprintf(
		"user=%s dbname=%s password=%s port=%s",
		os.Getenv("DBUSER"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPASS"),
		os.Getenv("DBPORT"),
	)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("CREATE EXTENSION postGIS")
	if err != nil {
		fmt.Println(err)
	}

	_, err = db.Exec("CREATE EXTENSION pgRouting")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Database extensions created.")

	cmd = exec.Command("osm2pgrouting", "-c", tagCostConf, "-p", os.Getenv("DBPORT"), "-d", os.Getenv("DBNAME"), "-f", "/tmp/filtered.osm", "-U", os.Getenv("DBUSER"), "-W", os.Getenv("DBPASS"))
	_, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Database filled.")

	return nil
}
