package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vishnupahwa/lakctl/cmd/commands"
	"os"
)

func main() {
	// rootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   "lakctl",
		Short: "Control applications",
		Long:  `Testing utility to control applications via HTTP.`,
	}
	commands.InitConfig()
	commands.AddCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
