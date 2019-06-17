package commands

import (
	"github.com/spf13/cobra"
	"github.com/thanosKontos/gravelmap/routing_prepare"
	"log"
	"os"
)

// createRoutingDBCommand defines the create route command.
func createRoutingDBCommand() *cobra.Command {
	createRoutingDBCmd := &cobra.Command{
		Use:   "create-db",
		Short: "create database",
		Long:  "create database",
	}

	createRoutingDBCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return createRoutingDBCmdRun()
	}

	return createRoutingDBCmd
}

// createRoutingDBCmdRun defines the command run actions.
func createRoutingDBCmdRun() error {
	gmPreparer, err := routing_prepare.NewGravelmapPreparer(
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
		os.Getenv("DBPORT"),
		os.Getenv("DBDEFAULTDBNAME"),
		logger,
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

	return nil
}
