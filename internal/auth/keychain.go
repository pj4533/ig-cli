package auth

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const serviceName = "ig-cli"

// KeychainStore defines the interface for secure credential storage.
type KeychainStore interface {
	Set(key, value string) error
	Get(key string) (string, error)
	Delete(key string) error
}

// OSKeychain implements KeychainStore using the OS keychain.
type OSKeychain struct{}

// Set stores a value in the OS keychain.
func (k *OSKeychain) Set(key, value string) error {
	return keyring.Set(serviceName, key, value)
}

// Get retrieves a value from the OS keychain.
func (k *OSKeychain) Get(key string) (string, error) {
	val, err := keyring.Get(serviceName, key)
	if err != nil {
		return "", fmt.Errorf("keychain get %q: %w", key, err)
	}
	return val, nil
}

// Delete removes a value from the OS keychain.
func (k *OSKeychain) Delete(key string) error {
	return keyring.Delete(serviceName, key)
}

// MockKeychain implements KeychainStore in-memory for testing.
type MockKeychain struct {
	store map[string]string
}

// NewMockKeychain creates a new MockKeychain.
func NewMockKeychain() *MockKeychain {
	return &MockKeychain{store: make(map[string]string)}
}

// Set stores a value in memory.
func (m *MockKeychain) Set(key, value string) error {
	m.store[key] = value
	return nil
}

// Get retrieves a value from memory.
func (m *MockKeychain) Get(key string) (string, error) {
	val, ok := m.store[key]
	if !ok {
		return "", fmt.Errorf("keychain get %q: not found", key)
	}
	return val, nil
}

// Delete removes a value from memory.
func (m *MockKeychain) Delete(key string) error {
	delete(m.store, key)
	return nil
}

// TokenKey returns the keychain key for a user's access token.
func TokenKey(username string) string {
	return fmt.Sprintf("token:%s", username)
}

// UserIDKey returns the keychain key for a user's Instagram user ID.
func UserIDKey(username string) string {
	return fmt.Sprintf("user_id:%s", username)
}

// TokenExpiryKey returns the keychain key for a token's expiry time.
func TokenExpiryKey(username string) string {
	return fmt.Sprintf("token_expiry:%s", username)
}
