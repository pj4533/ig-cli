package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var discoverCmd = &cobra.Command{
	Use:   "discover <username>",
	Short: "Look up a public Business/Creator account",
	Long:  "Discover public information about an Instagram Business or Creator account.",
	Args:  cobra.ExactArgs(1),
	RunE:  runDiscover,
}

func init() {
	rootCmd.AddCommand(discoverCmd)
}

func runDiscover(cmd *cobra.Command, args []string) error {
	client, userID, err := getClient()
	if err != nil {
		return err
	}

	discovery, err := client.DiscoverUser(userID, args[0])
	if err != nil {
		return fmt.Errorf("discovering user: %w", err)
	}

	return outputJSON(discovery)
}
