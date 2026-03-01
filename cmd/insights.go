package cmd

import "github.com/spf13/cobra"

var insightsCmd = &cobra.Command{
	Use:   "insights",
	Short: "View account insights and demographics",
	Long:  "Commands for viewing account-level insights and audience demographics.",
}

func init() {
	rootCmd.AddCommand(insightsCmd)
}
