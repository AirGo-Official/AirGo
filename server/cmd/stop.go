package cmd

import (
	"github.com/ppoonk/AirGo/utils/os_plugin"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop AirGo",
	Run: func(cmd *cobra.Command, args []string) {
		os_plugin.StopProcess("AirGo")
	},
}
