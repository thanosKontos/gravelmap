package commands

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Execute entry point for commands.
func Execute() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading .env file")
		os.Exit(1)
	}

	rootCommand := rootCommand()
	rootCommand.AddCommand(versionCommand())
	rootCommand.AddCommand(createRoutingDataCommand())
	rootCommand.AddCommand(createServerCommand())
	rootCommand.AddCommand(importElevationCommand())
	rootCommand.AddCommand(createGradeWaysCommand())
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// setUpCommand initializes and returns a new command object.
func rootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gravelmap",
		Short: "gravelmap is a route engine made as a composite of other sevices (osmium, postgis and pgrouting)",
		Long:  "gravelmap is a route engine made as a composite of other sevices (osmium, postgis and pgrouting)",
	}
}
