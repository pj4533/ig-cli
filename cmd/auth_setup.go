package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pj4533/ig-cli/internal/config"
	"github.com/spf13/cobra"
)

var authSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure Meta App ID and Secret",
	Long:  "Set up your Meta Developer App credentials for Instagram API access.",
	RunE:  runAuthSetup,
}

func init() {
	authCmd.AddCommand(authSetupCmd)
}

func runAuthSetup(cmd *cobra.Command, args []string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your Meta App ID: ")
	appID, _ := reader.ReadString('\n')
	appID = strings.TrimSpace(appID)

	if appID == "" {
		return fmt.Errorf("App ID cannot be empty")
	}

	fmt.Print("Enter your Meta App Secret: ")
	appSecret, _ := reader.ReadString('\n')
	appSecret = strings.TrimSpace(appSecret)

	if appSecret == "" {
		return fmt.Errorf("App Secret cannot be empty")
	}

	// Store App ID in config (not sensitive)
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	cfg.AppID = appID
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	// Store App Secret in keychain (sensitive)
	keychain := keychainFactory()
	if err := keychain.Set("app_secret", appSecret); err != nil {
		return fmt.Errorf("storing app secret: %w", err)
	}

	fmt.Println("Setup complete! Run 'ig auth add' to connect an Instagram account.")
	return nil
}
