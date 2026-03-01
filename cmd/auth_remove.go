package cmd

import (
	"fmt"

	"github.com/pj4533/ig-cli/internal/auth"
	"github.com/pj4533/ig-cli/internal/config"
	"github.com/spf13/cobra"
)

var authRemoveCmd = &cobra.Command{
	Use:   "remove <username>",
	Short: "Disconnect an Instagram account",
	Long:  "Remove an Instagram account and delete its stored credentials.",
	Args:  cobra.ExactArgs(1),
	RunE:  runAuthRemove,
}

func init() {
	authCmd.AddCommand(authRemoveCmd)
}

func runAuthRemove(cmd *cobra.Command, args []string) error {
	username := args[0]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if !cfg.RemoveAccount(username) {
		return fmt.Errorf("account %q not found", username)
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	// Clean up keychain entries (ignore errors for missing keys)
	keychain := keychainFactory()
	_ = keychain.Delete(auth.TokenKey(username))
	_ = keychain.Delete(auth.UserIDKey(username))
	_ = keychain.Delete(auth.TokenExpiryKey(username))

	fmt.Printf("Account %q removed.\n", username)
	return nil
}
