package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func setupTestConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	viper.Reset()
	return dir
}

func TestLoadEmptyConfig(t *testing.T) {
	setupTestConfig(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.AppID != "" {
		t.Errorf("AppID = %q, want empty", cfg.AppID)
	}
	if len(cfg.Accounts) != 0 {
		t.Errorf("Accounts length = %d, want 0", len(cfg.Accounts))
	}
}

func TestSaveAndLoad(t *testing.T) {
	setupTestConfig(t)

	cfg := &Config{
		AppID:          "test-app-id",
		DefaultAccount: "testuser",
		Accounts: []Account{
			{Username: "testuser", UserID: "12345"},
		},
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	// Reset viper to force a fresh read
	viper.Reset()

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if loaded.AppID != cfg.AppID {
		t.Errorf("AppID = %q, want %q", loaded.AppID, cfg.AppID)
	}
	if loaded.DefaultAccount != cfg.DefaultAccount {
		t.Errorf("DefaultAccount = %q, want %q", loaded.DefaultAccount, cfg.DefaultAccount)
	}
	if len(loaded.Accounts) != 1 {
		t.Fatalf("Accounts length = %d, want 1", len(loaded.Accounts))
	}
	if loaded.Accounts[0].Username != "testuser" {
		t.Errorf("Username = %q, want %q", loaded.Accounts[0].Username, "testuser")
	}
}

func TestAddAccount(t *testing.T) {
	cfg := &Config{}

	cfg.AddAccount("user1", "111")
	if len(cfg.Accounts) != 1 {
		t.Fatalf("Accounts length = %d, want 1", len(cfg.Accounts))
	}
	if cfg.DefaultAccount != "user1" {
		t.Errorf("DefaultAccount = %q, want %q", cfg.DefaultAccount, "user1")
	}

	// Add second account - default should remain user1
	cfg.AddAccount("user2", "222")
	if len(cfg.Accounts) != 2 {
		t.Fatalf("Accounts length = %d, want 2", len(cfg.Accounts))
	}
	if cfg.DefaultAccount != "user1" {
		t.Errorf("DefaultAccount = %q, want %q", cfg.DefaultAccount, "user1")
	}

	// Adding existing account should update, not duplicate
	cfg.AddAccount("user1", "111-updated")
	if len(cfg.Accounts) != 2 {
		t.Fatalf("Accounts length = %d, want 2 (should not duplicate)", len(cfg.Accounts))
	}
	if cfg.Accounts[0].UserID != "111-updated" {
		t.Errorf("UserID = %q, want %q", cfg.Accounts[0].UserID, "111-updated")
	}
}

func TestRemoveAccount(t *testing.T) {
	cfg := &Config{
		DefaultAccount: "user1",
		Accounts: []Account{
			{Username: "user1", UserID: "111"},
			{Username: "user2", UserID: "222"},
		},
	}

	// Remove default account - should switch default to next
	if !cfg.RemoveAccount("user1") {
		t.Error("RemoveAccount returned false, want true")
	}
	if len(cfg.Accounts) != 1 {
		t.Fatalf("Accounts length = %d, want 1", len(cfg.Accounts))
	}
	if cfg.DefaultAccount != "user2" {
		t.Errorf("DefaultAccount = %q, want %q", cfg.DefaultAccount, "user2")
	}

	// Remove last account
	if !cfg.RemoveAccount("user2") {
		t.Error("RemoveAccount returned false, want true")
	}
	if cfg.DefaultAccount != "" {
		t.Errorf("DefaultAccount = %q, want empty", cfg.DefaultAccount)
	}

	// Remove non-existent
	if cfg.RemoveAccount("nobody") {
		t.Error("RemoveAccount returned true for non-existent, want false")
	}
}

func TestGetAccount(t *testing.T) {
	cfg := &Config{
		Accounts: []Account{
			{Username: "user1", UserID: "111"},
		},
	}

	acct := cfg.GetAccount("user1")
	if acct == nil {
		t.Fatal("GetAccount returned nil, want account")
	}
	if acct.UserID != "111" {
		t.Errorf("UserID = %q, want %q", acct.UserID, "111")
	}

	if cfg.GetAccount("nobody") != nil {
		t.Error("GetAccount returned non-nil for non-existent account")
	}
}

func TestActiveAccount(t *testing.T) {
	cfg := &Config{
		DefaultAccount: "user1",
		Accounts: []Account{
			{Username: "user1", UserID: "111"},
			{Username: "user2", UserID: "222"},
		},
	}

	// Default
	acct, err := cfg.ActiveAccount("")
	if err != nil {
		t.Fatalf("ActiveAccount error: %v", err)
	}
	if acct.Username != "user1" {
		t.Errorf("Username = %q, want %q", acct.Username, "user1")
	}

	// Override
	acct, err = cfg.ActiveAccount("user2")
	if err != nil {
		t.Fatalf("ActiveAccount error: %v", err)
	}
	if acct.Username != "user2" {
		t.Errorf("Username = %q, want %q", acct.Username, "user2")
	}

	// Non-existent override
	_, err = cfg.ActiveAccount("nobody")
	if err == nil {
		t.Error("ActiveAccount should error for non-existent account")
	}

	// No accounts configured
	empty := &Config{}
	_, err = empty.ActiveAccount("")
	if err == nil {
		t.Error("ActiveAccount should error when no accounts configured")
	}
}

func TestSaveCreatesDir(t *testing.T) {
	dir := setupTestConfig(t)

	cfg := &Config{AppID: "test"}
	if err := Save(cfg); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	configDir := filepath.Join(dir, ".ig-cli")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Error("config directory was not created")
	}
}
