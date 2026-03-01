package cmd

import (
	"fmt"

	"github.com/pj4533/ig-cli/internal/auth"
	"github.com/pj4533/ig-cli/internal/config"
	"github.com/spf13/cobra"
)

var authAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Connect an Instagram account via OAuth",
	Long:  "Start the OAuth flow to connect a new Instagram Business or Creator account.",
	RunE:  runAuthAdd,
}

func init() {
	authCmd.AddCommand(authAddCmd)
}

func runAuthAdd(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if cfg.AppID == "" {
		return fmt.Errorf("no App ID configured; run 'ig auth setup' first")
	}

	keychain := keychainFactory()
	appSecret, err := keychain.Get("app_secret")
	if err != nil {
		return fmt.Errorf("no App Secret found; run 'ig auth setup' first")
	}

	client := clientFactory("")
	flow := &auth.OAuthFlow{
		AppID:     cfg.AppID,
		AppSecret: appSecret,
		Client:    client,
		Keychain:  keychain,
	}

	result, err := oauthFlowRunner(flow)
	if err != nil {
		return fmt.Errorf("OAuth flow failed: %w", err)
	}

	// Store token in keychain
	tm := auth.NewTokenManager(keychain, client)
	if err := tm.StoreToken(result.Username, result.Token, result.ExpiresIn); err != nil {
		return fmt.Errorf("storing token: %w", err)
	}

	// Store user ID in keychain
	if err := keychain.Set(auth.UserIDKey(result.Username), result.UserID); err != nil {
		return fmt.Errorf("storing user ID: %w", err)
	}

	// Add account to config
	cfg.AddAccount(result.Username, result.UserID)
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Printf("Successfully connected account: %s\n", result.Username)
	return nil
}
