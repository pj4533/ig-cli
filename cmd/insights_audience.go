package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var insightsAudienceCmd = &cobra.Command{
	Use:   "audience",
	Short: "Get audience demographics",
	Long:  "View audience demographic data including age, gender, city, and country breakdowns.",
	RunE:  runInsightsAudience,
}

func init() {
	insightsCmd.AddCommand(insightsAudienceCmd)
}

func runInsightsAudience(cmd *cobra.Command, args []string) error {
	client, userID, err := getClient()
	if err != nil {
		return err
	}

	demographics, err := client.GetAudienceDemographics(userID)
	if err != nil {
		return fmt.Errorf("getting audience demographics: %w", err)
	}

	return outputJSON(demographics)
}
