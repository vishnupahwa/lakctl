package commands

import (
	"os"

	"github.com/spf13/cobra"
)

func addCompletion(topLevel *cobra.Command) {
	var zsh bool
	// completionCmd represents the completion command
	var completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if zsh {
				return topLevel.GenZshCompletion(os.Stdout)
			}
			return topLevel.GenBashCompletion(os.Stdout)
		},
	}
	completionCmd.Flags().BoolVar(&zsh, "zsh", false, "Output zsh completion instead")
	topLevel.AddCommand(completionCmd)
}
