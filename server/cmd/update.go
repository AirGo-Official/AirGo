package cmd

import (
	"github.com/ppoonk/AirGo/app"
	"github.com/spf13/cobra"
)

func init() {
	updateCmd.Flags().StringVar(&startConfigPath, "config", "config.yaml", "config.yaml directory to read from")
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update AirGo",
	Run: func(cmd *cobra.Command, args []string) {
		update()
	},
}

func update() {

	newApp := app.NewApp()

	newApp.InitConfig(startConfigPath)

	newApp.InitLogrus()

	newApp.InitLocalCache()

	newApp.InitRouter()

	newApp.ConnectDatabase()

	newApp.Update()

}
