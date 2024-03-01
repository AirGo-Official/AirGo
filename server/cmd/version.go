package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

const v = "old-version"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version of AirGo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(v)
	},
}
