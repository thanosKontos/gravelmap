package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Execute entry point for commands.
func Execute() {
	rootCommand := setUpCommand()
	rootCommand.AddCommand(versionCommand())
	rootCommand.AddCommand(mapDrawerCommand())
	if err := rootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// setUpCommand initializes and returns a new command object.
func setUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "gravelmap",
		Short: "gravelmap is a routing engine made as a composite of other sevices (osmium, postgis and pgrouting)",
		Long: "gravelmap is a routing engine made as a composite of other sevices (osmium, postgis and pgrouting)",
	}
}
