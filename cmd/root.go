package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCommand = &cobra.Command{
		Use:   "simpleAPI",
		Short: "simpleAPI service",
	}
	configPath string
)

func init() {
	rootCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "app.env", "config file path")

	rootCommand.AddCommand(serveCommand)
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		fmt.Printf("failed to execute root command: %s\n", err.Error())
		os.Exit(1)
	}
}
