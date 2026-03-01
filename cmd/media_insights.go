package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var mediaInsightsCmd = &cobra.Command{
	Use:   "insights <media-id>",
	Short: "Get insights for a media post",
	Long:  "View detailed metrics for a specific Instagram media post.",
	Args:  cobra.ExactArgs(1),
	RunE:  runMediaInsights,
}

func init() {
	mediaCmd.AddCommand(mediaInsightsCmd)
}

func runMediaInsights(cmd *cobra.Command, args []string) error {
	client, _, err := getClient()
	if err != nil {
		return err
	}

	insights, err := client.GetMediaInsights(args[0])
	if err != nil {
		return fmt.Errorf("getting media insights: %w", err)
	}

	return outputJSON(insights)
}
