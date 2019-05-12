package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

// versionCommand defines the version command.
func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "the version of gravelmap",
		Long:  "the version of gravelmap",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version: ", "0.0.2")
		},
	}
}
