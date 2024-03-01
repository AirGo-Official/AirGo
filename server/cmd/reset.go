package cmd

import (
	"github.com/ppoonk/AirGo/initialize"
	"github.com/ppoonk/AirGo/service/admin_logic"
	"github.com/spf13/cobra"
)

func init() {
	resetCmd.Flags().StringVar(&startConfigPath, "config", "config.yaml", "config.yaml directory to read from")
	resetCmd.Flags().Bool("resetAdmin", false, "reset administrator account and password")
	resetCmd.MarkFlagsOneRequired("resetAdmin")
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "reset resources",
	Run: func(cmd *cobra.Command, args []string) {
		var adminUserService admin_logic.User
		initialize.InitializeDB(startConfigPath)
		adminUserService.ResetAdminPassword()
	},
	//Args: cobra.MaximumNArgs(1),
}
