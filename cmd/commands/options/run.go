package options

import "github.com/spf13/cobra"

type Run struct {
	Command string
}

func AddRunArg(cmd *cobra.Command, r *Run) {
	cmd.Flags().StringVarP(&r.Command, "run", "r", "", "Run command to start")
	_ = cmd.MarkFlagRequired("run")
}
