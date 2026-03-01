package cmd

import "github.com/spf13/cobra"

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage Instagram authentication",
	Long:  "Commands for setting up and managing Instagram account authentication.",
}

func init() {
	rootCmd.AddCommand(authCmd)
}
