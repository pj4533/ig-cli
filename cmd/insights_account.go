package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var insightsAccountPeriod string

var insightsAccountCmd = &cobra.Command{
	Use:   "account",
	Short: "Get account-level insights",
	Long:  "View account-level metrics including reach, views, and follower growth.",
	RunE:  runInsightsAccount,
}

func init() {
	insightsAccountCmd.Flags().StringVar(&insightsAccountPeriod, "period", "day", "Time period: day, week, or days_28")
	insightsCmd.AddCommand(insightsAccountCmd)
}

func runInsightsAccount(cmd *cobra.Command, args []string) error {
	client, userID, err := getClient()
	if err != nil {
		return err
	}

	insights, err := client.GetAccountInsights(userID, insightsAccountPeriod)
	if err != nil {
		return fmt.Errorf("getting account insights: %w", err)
	}

	return outputJSON(insights)
}
