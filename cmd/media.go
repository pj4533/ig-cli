package cmd

import "github.com/spf13/cobra"

var mediaCmd = &cobra.Command{
	Use:   "media",
	Short: "Manage and view media posts",
	Long:  "Commands for listing media posts and viewing their insights.",
}

func init() {
	rootCmd.AddCommand(mediaCmd)
}
