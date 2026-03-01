package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pj4533/ig-cli/internal/auth"
	"github.com/pj4533/ig-cli/internal/config"
	"github.com/spf13/cobra"
)

var authListCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected Instagram accounts",
	Long:  "Show all connected Instagram accounts with token expiry information.",
	RunE:  runAuthList,
}

func init() {
	authCmd.AddCommand(authListCmd)
}

func runAuthList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if len(cfg.Accounts) == 0 {
		fmt.Println("No accounts connected. Run 'ig auth add' to connect one.")
		return nil
	}

	keychain := keychainFactory()

	type accountInfo struct {
		Username  string `json:"username"`
		UserID    string `json:"user_id"`
		Default   bool   `json:"default"`
		ExpiresAt string `json:"expires_at,omitempty"`
		ExpiresIn string `json:"expires_in,omitempty"`
	}

	var accounts []accountInfo
	for _, acct := range cfg.Accounts {
		info := accountInfo{
			Username: acct.Username,
			UserID:   acct.UserID,
			Default:  acct.Username == cfg.DefaultAccount,
		}

		if expiryStr, err := keychain.Get(auth.TokenExpiryKey(acct.Username)); err == nil {
			if expiry, err := strconv.ParseInt(expiryStr, 10, 64); err == nil {
				expiryTime := time.Unix(expiry, 0)
				info.ExpiresAt = expiryTime.Format(time.RFC3339)
				remaining := time.Until(expiryTime)
				if remaining > 0 {
					info.ExpiresIn = fmt.Sprintf("%.0f days", remaining.Hours()/24)
				} else {
					info.ExpiresIn = "expired"
				}
			}
		}

		accounts = append(accounts, info)
	}

	return outputJSON(accounts)
}
