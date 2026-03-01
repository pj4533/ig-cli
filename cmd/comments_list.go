package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commentsListLimit int

var commentsListCmd = &cobra.Command{
	Use:   "list <media-id>",
	Short: "List comments on a media post",
	Long:  "List all comments on a specific Instagram media post.",
	Args:  cobra.ExactArgs(1),
	RunE:  runCommentsList,
}

func init() {
	commentsListCmd.Flags().IntVar(&commentsListLimit, "limit", 0, "Maximum number of comments to return")
	commentsCmd.AddCommand(commentsListCmd)
}

func runCommentsList(cmd *cobra.Command, args []string) error {
	client, _, err := getClient()
	if err != nil {
		return err
	}

	comments, err := client.ListComments(args[0], commentsListLimit)
	if err != nil {
		return fmt.Errorf("listing comments: %w", err)
	}

	return outputJSON(comments)
}
