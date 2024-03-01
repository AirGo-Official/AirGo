package cmd

import (
	"github.com/ppoonk/AirGo/initialize"
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().StringVar(&startConfigPath, "config", "config.yaml", "config.yaml directory to read from")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update AirGo",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.InitializeUpdate(startConfigPath)
	},
}
