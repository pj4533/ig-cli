package cmd

import "github.com/spf13/cobra"

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "View comments and replies",
	Long:  "Commands for listing comments on posts and viewing replies.",
}

func init() {
	rootCmd.AddCommand(commentsCmd)
}
