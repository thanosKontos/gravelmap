package commands

import (
	"fmt"
	"github.com/thanosKontos/gravelmap"
	"github.com/thanosKontos/gravelmap/log"

	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var logger gravelmap.Logger

// Execute entry point for commands.
func Execute() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading .env file")
		os.Exit(1)
	}

	rootCommand := rootCommand()
	rootCommand.AddCommand(versionCommand())
	rootCommand.AddCommand(createRoutingDBCommand())
	rootCommand.AddCommand(createServerCommand())
	rootCommand.AddCommand(createWebServerNewCommand())
	rootCommand.AddCommand(importElevationCommand())
	rootCommand.AddCommand(createGradeWaysCommand())
	rootCommand.AddCommand(importRoutingDataCommand())
	rootCommand.AddCommand(createFilterOSMCommand())
	rootCommand.AddCommand(dijkstraPocCommand())
	rootCommand.AddCommand(osmParserCommand())
	rootCommand.AddCommand(searchCommand())
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// setUpCommand initializes and returns a new command object.
func rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gravelmap",
		Short: "gravelmap is a route engine made as a composite of other sevices (osmium, postgis and pgrouting)",
		Long:  "gravelmap is a route engine made as a composite of other sevices (osmium, postgis and pgrouting)",
	}

	var verbose bool
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return persistentPreRunECommand(cmd.Name(), verbose)
	}

	return rootCmd
}

// persistentPreRunECommand defines the root command actions before the run command.
func persistentPreRunECommand(cmdName string, verbose bool) error {
	if verbose {
		logger = log.NewDebugCLI()
	} else {
		logger = log.NewNullCLI()
	}

	return nil
}
