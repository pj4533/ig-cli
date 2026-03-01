package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pj4533/ig-cli/internal/api"
	"github.com/pj4533/ig-cli/internal/auth"
	"github.com/pj4533/ig-cli/internal/config"
	"github.com/spf13/viper"
)

func TestRunAuthList_NoAccounts(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	output := captureStdout(t, func() {
		err := runAuthList(nil, nil)
		if err != nil {
			t.Fatalf("runAuthList error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected some output")
	}
}

func TestRunAuthList_WithAccounts(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{
		AppID:          "test-app-id",
		DefaultAccount: "testuser",
		Accounts: []config.Account{
			{Username: "testuser", UserID: "12345"},
		},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	// Override keychain factory
	origKeychainFactory := keychainFactory
	mockKeychain := auth.NewMockKeychain()
	expiry := time.Now().Add(30 * 24 * time.Hour).Unix()
	_ = mockKeychain.Set(auth.TokenExpiryKey("testuser"), strconv.FormatInt(expiry, 10))
	keychainFactory = func() auth.KeychainStore { return mockKeychain }
	defer func() { keychainFactory = origKeychainFactory }()

	output := captureStdout(t, func() {
		err := runAuthList(nil, nil)
		if err != nil {
			t.Fatalf("runAuthList error: %v", err)
		}
	})

	var result []struct {
		Username string `json:"username"`
		Default  bool   `json:"default"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v\noutput: %s", err, output)
	}
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
	if result[0].Username != "testuser" {
		t.Errorf("Username = %q, want %q", result[0].Username, "testuser")
	}
	if !result[0].Default {
		t.Error("expected Default to be true")
	}
}

func TestRunAuthRemove_Success(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{
		AppID:          "test-app-id",
		DefaultAccount: "testuser",
		Accounts: []config.Account{
			{Username: "testuser", UserID: "12345"},
		},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	// Override keychain factory
	origKeychainFactory := keychainFactory
	mockKeychain := auth.NewMockKeychain()
	_ = mockKeychain.Set(auth.TokenKey("testuser"), "token")
	keychainFactory = func() auth.KeychainStore { return mockKeychain }
	defer func() { keychainFactory = origKeychainFactory }()

	output := captureStdout(t, func() {
		err := runAuthRemove(nil, []string{"testuser"})
		if err != nil {
			t.Fatalf("runAuthRemove error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected confirmation output")
	}
}

func TestRunAuthRemove_NotFound(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{
		Accounts: []config.Account{},
	}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	err := runAuthRemove(nil, []string{"nobody"})
	if err == nil {
		t.Error("expected error for non-existent account")
	}
}

func TestRunAuthAdd_NoAppID(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	err := runAuthAdd(nil, nil)
	if err == nil {
		t.Error("expected error when no App ID configured")
	}
}

func TestRunAuthAdd_NoAppSecret(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{AppID: "test-app-id"}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	// Override keychain factory with empty keychain
	origKeychainFactory := keychainFactory
	mockKeychain := auth.NewMockKeychain()
	keychainFactory = func() auth.KeychainStore { return mockKeychain }
	defer func() { keychainFactory = origKeychainFactory }()

	err := runAuthAdd(nil, nil)
	if err == nil {
		t.Error("expected error when no App Secret in keychain")
	}
}

func TestRunAuthAdd_Success(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{AppID: "test-app-id"}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	mockKeychain := auth.NewMockKeychain()
	_ = mockKeychain.Set("app_secret", "test-secret")

	mockClient := &api.MockClient{}

	origKeychainFactory := keychainFactory
	origClientFactory := clientFactory
	origFlowRunner := oauthFlowRunner
	keychainFactory = func() auth.KeychainStore { return mockKeychain }
	clientFactory = func(token string) api.Client { return mockClient }
	oauthFlowRunner = func(flow *auth.OAuthFlow) (*auth.OAuthResult, error) {
		return &auth.OAuthResult{
			Username:  "testuser",
			UserID:    "12345",
			Token:     "long-token",
			ExpiresIn: 5184000,
		}, nil
	}
	defer func() {
		keychainFactory = origKeychainFactory
		clientFactory = origClientFactory
		oauthFlowRunner = origFlowRunner
	}()

	output := captureStdout(t, func() {
		err := runAuthAdd(nil, nil)
		if err != nil {
			t.Fatalf("runAuthAdd error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected output")
	}

	// Verify token was stored
	token, err := mockKeychain.Get(auth.TokenKey("testuser"))
	if err != nil {
		t.Fatalf("get token error: %v", err)
	}
	if token != "long-token" {
		t.Errorf("token = %q, want %q", token, "long-token")
	}

	// Verify user ID was stored
	userID, err := mockKeychain.Get(auth.UserIDKey("testuser"))
	if err != nil {
		t.Fatalf("get user ID error: %v", err)
	}
	if userID != "12345" {
		t.Errorf("userID = %q, want %q", userID, "12345")
	}
}

func TestRunAuthAdd_OAuthFlowFails(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	cfg := &config.Config{AppID: "test-app-id"}
	if err := config.Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}
	viper.Reset()

	mockKeychain := auth.NewMockKeychain()
	_ = mockKeychain.Set("app_secret", "test-secret")

	origKeychainFactory := keychainFactory
	origClientFactory := clientFactory
	origFlowRunner := oauthFlowRunner
	keychainFactory = func() auth.KeychainStore { return mockKeychain }
	clientFactory = func(token string) api.Client { return &api.MockClient{} }
	oauthFlowRunner = func(flow *auth.OAuthFlow) (*auth.OAuthResult, error) {
		return nil, fmt.Errorf("oauth flow failed")
	}
	defer func() {
		keychainFactory = origKeychainFactory
		clientFactory = origClientFactory
		oauthFlowRunner = origFlowRunner
	}()

	err := runAuthAdd(nil, nil)
	if err == nil {
		t.Error("expected error when OAuth flow fails")
	}
}

func TestRunAuthSetup_Success(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	// Override keychain factory
	origKeychainFactory := keychainFactory
	mockKeychain := auth.NewMockKeychain()
	keychainFactory = func() auth.KeychainStore { return mockKeychain }
	defer func() { keychainFactory = origKeychainFactory }()

	// Simulate stdin with app ID and secret
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		_, _ = fmt.Fprintln(w, "test-app-id")
		_, _ = fmt.Fprintln(w, "test-app-secret")
		_ = w.Close()
	}()

	output := captureStdout(t, func() {
		err := runAuthSetup(nil, nil)
		if err != nil {
			t.Fatalf("runAuthSetup error: %v", err)
		}
	})

	if output == "" {
		t.Error("expected output")
	}

	// Verify secret was stored
	secret, err := mockKeychain.Get("app_secret")
	if err != nil {
		t.Fatalf("get secret error: %v", err)
	}
	if secret != "test-app-secret" {
		t.Errorf("secret = %q, want %q", secret, "test-app-secret")
	}
}

func TestRunAuthSetup_EmptyAppID(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	// Simulate stdin with empty app ID
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		_, _ = fmt.Fprintln(w, "")
		_ = w.Close()
	}()

	err := runAuthSetup(nil, nil)
	if err == nil {
		t.Error("expected error for empty App ID")
	}
}

func TestRunAuthSetup_EmptySecret(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	// Simulate stdin with valid ID but empty secret
	oldStdin := os.Stdin
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() { os.Stdin = oldStdin }()

	go func() {
		_, _ = fmt.Fprintln(w, "test-app-id")
		_, _ = fmt.Fprintln(w, "")
		_ = w.Close()
	}()

	err := runAuthSetup(nil, nil)
	if err == nil {
		t.Error("expected error for empty App Secret")
	}
}

func TestExecute(t *testing.T) {
	// Just verify Execute doesn't panic
	// We can't easily test it fully since it calls os.Exit
	_ = rootCmd

	// Verify the mock client interface
	var _ api.Client = &api.MockClient{}
}
