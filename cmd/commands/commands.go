package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddCommands(topLevel *cobra.Command) {
	addStart(topLevel)
}

// initConfig reads in config file and ENV variables if set.
func InitConfig() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./")
	viper.SetConfigName(".lakctl")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
