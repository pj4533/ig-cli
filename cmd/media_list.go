package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mediaListLimit int

var mediaListCmd = &cobra.Command{
	Use:   "list",
	Short: "List media posts",
	Long:  "List Instagram media posts with like and comment counts.",
	RunE:  runMediaList,
}

func init() {
	mediaListCmd.Flags().IntVar(&mediaListLimit, "limit", 0, "Maximum number of posts to return")
	mediaCmd.AddCommand(mediaListCmd)
}

func runMediaList(cmd *cobra.Command, args []string) error {
	client, userID, err := getClient()
	if err != nil {
		return err
	}

	media, err := client.ListMedia(userID, mediaListLimit)
	if err != nil {
		return fmt.Errorf("listing media: %w", err)
	}

	return outputJSON(media)
}
