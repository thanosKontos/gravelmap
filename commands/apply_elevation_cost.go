package commands

import (
	"github.com/spf13/cobra"
)

// applyElevationCostCommand defines the create routing command.
func applyElevationCostCommand() *cobra.Command {
	applyElevationCostCmd := &cobra.Command{
		Use:   "apply-elevation-cost",
		Short: "add elevation cost to routing database",
		Long:  "add elevation cost to routing database",
	}

	applyElevationCostCmd.RunE = func(cmd *cobra.Command, args []string) error {
		return applyElevationCostCmdRun()
	}

	return applyElevationCostCmd
}

// applyElevationCostCmdRun defines the command run actions.
func applyElevationCostCmdRun() error {
	// TBD

	return nil
}
