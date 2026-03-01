package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var commentsRepliesLimit int

var commentsRepliesCmd = &cobra.Command{
	Use:   "replies <comment-id>",
	Short: "List replies to a comment",
	Long:  "List all replies to a specific comment.",
	Args:  cobra.ExactArgs(1),
	RunE:  runCommentsReplies,
}

func init() {
	commentsRepliesCmd.Flags().IntVar(&commentsRepliesLimit, "limit", 0, "Maximum number of replies to return")
	commentsCmd.AddCommand(commentsRepliesCmd)
}

func runCommentsReplies(cmd *cobra.Command, args []string) error {
	client, _, err := getClient()
	if err != nil {
		return err
	}

	replies, err := client.ListReplies(args[0], commentsRepliesLimit)
	if err != nil {
		return fmt.Errorf("listing replies: %w", err)
	}

	return outputJSON(replies)
}
