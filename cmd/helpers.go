package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/pj4533/ig-cli/internal/api"
	"github.com/pj4533/ig-cli/internal/auth"
	"github.com/pj4533/ig-cli/internal/config"
)

// clientFactory allows overriding the client creation for testing.
var clientFactory func(token string) api.Client

// keychainFactory allows overriding keychain for testing.
var keychainFactory func() auth.KeychainStore

func init() {
	clientFactory = func(token string) api.Client {
		return api.NewGraphClient(token)
	}
	keychainFactory = func() auth.KeychainStore {
		return &auth.OSKeychain{}
	}
}

// getClient resolves the active account and returns a configured API client and user ID.
func getClient() (api.Client, string, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, "", fmt.Errorf("loading config: %w", err)
	}

	acct, err := cfg.ActiveAccount(accountFlag)
	if err != nil {
		return nil, "", err
	}

	slog.Debug("Using account", "username", acct.Username, "user_id", acct.UserID)

	keychain := keychainFactory()
	tempClient := clientFactory("")
	tm := auth.NewTokenManager(keychain, tempClient)

	token, err := tm.GetValidToken(acct.Username)
	if err != nil {
		return nil, "", fmt.Errorf("getting token for %q: %w", acct.Username, err)
	}

	client := clientFactory(token)
	return client, acct.UserID, nil
}

// outputJSON prints data as pretty-printed JSON to stdout.
func outputJSON(data interface{}) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// outputError prints an error as JSON to stderr.
func outputError(err error) {
	errObj := map[string]string{"error": err.Error()}
	data, _ := json.MarshalIndent(errObj, "", "  ")
	fmt.Fprintln(os.Stderr, string(data))
}
