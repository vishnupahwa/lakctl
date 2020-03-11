package options

import "github.com/spf13/cobra"

// Run struct contains options regarding the main command to be controlled.
type Run struct {
	Command string
}

func AddRunArg(cmd *cobra.Command, r *Run) {
	cmd.Flags().StringVarP(&r.Command, "run", "r", "", "Run command to start")
	_ = cmd.MarkFlagRequired("run")
}
