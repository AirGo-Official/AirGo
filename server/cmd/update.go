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
		// 开发时，通过命令行升级数据库 role_and_menu 、 menu 以及 casbin_rule。因为开发时，经常修改api接口和菜单。
		initialize.InitializeUpdate(startConfigPath)
	},
}
