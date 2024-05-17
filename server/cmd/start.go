package cmd

import (
	"fmt"
	"github.com/ppoonk/AirGo/app"
	"github.com/ppoonk/AirGo/constant"
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
		fmt.Println("[AirGo version] " + constant.V)
		start()
	},
}

func start() {

	newApp := app.NewApp()

	newApp.InitConfig(startConfigPath)

	newApp.InitLogrus()

	newApp.InitLocalCache()

	newApp.InitRouter()

	newApp.ConnectDatabase()

	newApp.InitGlobalVariable()

	newApp.InitTasks()

	newApp.Start()
}
