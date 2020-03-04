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
	rootCommand.AddCommand(createWebServerCommand())
	rootCommand.AddCommand(createFilterOSMCommand())
	rootCommand.AddCommand(importRoutingDataCommand())
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// setUpCommand initializes and returns a new command object.
func rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "gravelmap",
		Short: "gravelmap is a route engine",
		Long:  "gravelmap is a route engine made for off-road adventurers (mountain bikers, SUV vehicles, hikers)",
	}

	var verboseLevel string
	rootCmd.PersistentFlags().StringVarP(&verboseLevel, "verbose-level", "v", "error", "verbose level")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return persistentPreRunECommand(cmd.Name(), verboseLevel)
	}

	return rootCmd
}

// persistentPreRunECommand defines the root command actions before the run command.
func persistentPreRunECommand(cmdName string, verboseLevel string) error {
	logger = log.NewStdout(verboseLevel)

	return nil
}
