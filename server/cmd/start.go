package cmd

import (
	"github.com/ppoonk/AirGo/initialize"
	"github.com/spf13/cobra"
)

func init() {
	startCmd.Flags().StringVar(&startConfigPath, "config", "config.yaml", "config.yaml directory to read from")
}

var (
	startConfigPath string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start AirGo",
	Run: func(cmd *cobra.Command, args []string) {
		initialize.InitializeAll(startConfigPath) //初始化系统资源并启动路由
	},
}
