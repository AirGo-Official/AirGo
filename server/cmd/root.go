package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(resetCmd)
}

var rootCmd = &cobra.Command{
	Use:   "AirGo",
	Short: "AirGo is a modern multi-user agent panel.",
	Long: `[ AirGo is a modern multi-user agent panel.]
[ Simple and easy to operate.]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args:", args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
