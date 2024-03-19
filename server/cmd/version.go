package cmd

import (
	"fmt"
	"github.com/ppoonk/AirGo/constant"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version of AirGo",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(constant.V)
	},
}
