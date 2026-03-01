package auth

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/pj4533/ig-cli/internal/api"
)

// TokenManager handles token retrieval and auto-refresh.
type TokenManager struct {
	keychain KeychainStore
	client   api.Client
}

// NewTokenManager creates a new TokenManager.
func NewTokenManager(keychain KeychainStore, client api.Client) *TokenManager {
	return &TokenManager{
		keychain: keychain,
		client:   client,
	}
}

// GetValidToken returns a valid access token for the given username,
// refreshing it if it's within 24 hours of expiry.
func (tm *TokenManager) GetValidToken(username string) (string, error) {
	token, err := tm.keychain.Get(TokenKey(username))
	if err != nil {
		return "", fmt.Errorf("no token found for %q: %w", username, err)
	}

	// Check expiry
	expiryStr, err := tm.keychain.Get(TokenExpiryKey(username))
	if err != nil {
		slog.Debug("No expiry found, using token as-is", "username", username)
		return token, nil
	}

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		slog.Debug("Invalid expiry format, using token as-is", "username", username)
		return token, nil
	}

	expiryTime := time.Unix(expiry, 0)
	timeUntilExpiry := time.Until(expiryTime)

	slog.Debug("Token expiry check",
		"username", username,
		"expires", expiryTime.Format(time.RFC3339),
		"hours_remaining", timeUntilExpiry.Hours(),
	)

	// Refresh if within 24 hours of expiry
	if timeUntilExpiry < 24*time.Hour {
		slog.Info("Token near expiry, refreshing", "username", username)
		refreshed, err := tm.client.RefreshLongLivedToken(token)
		if err != nil {
			slog.Warn("Token refresh failed, using existing token", "error", err)
			return token, nil
		}

		if err := tm.StoreToken(username, refreshed.AccessToken, refreshed.ExpiresIn); err != nil {
			slog.Warn("Failed to store refreshed token", "error", err)
		}

		return refreshed.AccessToken, nil
	}

	return token, nil
}

// StoreToken saves a token and its expiry to the keychain.
func (tm *TokenManager) StoreToken(username, token string, expiresIn int64) error {
	if err := tm.keychain.Set(TokenKey(username), token); err != nil {
		return fmt.Errorf("storing token: %w", err)
	}

	expiry := time.Now().Add(time.Duration(expiresIn) * time.Second).Unix()
	if err := tm.keychain.Set(TokenExpiryKey(username), strconv.FormatInt(expiry, 10)); err != nil {
		return fmt.Errorf("storing token expiry: %w", err)
	}

	return nil
}
