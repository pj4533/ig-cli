package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pj4533/ig-cli/internal/api"
	"github.com/pj4533/ig-cli/internal/models"
)

// noopBrowser is a no-op browser opener for tests.
func noopBrowser(url string) error {
	return nil
}

func TestOAuthFlow_Run_Success(t *testing.T) {
	keychain := NewMockKeychain()

	mockClient := &api.MockClient{
		ExchangeCodeForTokenFn: func(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error) {
			if code != "test-code" {
				t.Errorf("code = %q, want %q", code, "test-code")
			}
			return &models.TokenResponse{
				AccessToken: "short-token",
				ExpiresIn:   3600,
			}, nil
		},
		ExchangeForLongLivedTokenFn: func(appID, appSecret, shortToken string) (*models.TokenResponse, error) {
			if shortToken != "short-token" {
				t.Errorf("shortToken = %q, want %q", shortToken, "short-token")
			}
			return &models.TokenResponse{
				AccessToken: "long-token",
				ExpiresIn:   5184000,
			}, nil
		},
		GetUserProfileFn: func(token string) (*models.UserProfile, error) {
			if token != "long-token" {
				t.Errorf("token = %q, want %q", token, "long-token")
			}
			return &models.UserProfile{
				ID:       "12345",
				Username: "testuser",
				Name:     "Test User",
			}, nil
		},
	}

	flow := &OAuthFlow{
		AppID:       "test-app-id",
		AppSecret:   "test-app-secret",
		Client:      mockClient,
		Keychain:    keychain,
		OpenBrowser: noopBrowser,
	}

	// Run the flow in a goroutine
	resultChan := make(chan *OAuthResult, 1)
	errChan := make(chan error, 1)

	go func() {
		result, err := flow.Run()
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- result
	}()

	// Give the server a moment to start
	time.Sleep(200 * time.Millisecond)

	// Simulate the OAuth callback
	resp, err := http.Get("http://localhost:8080/callback?code=test-code")
	if err != nil {
		t.Fatalf("callback request error: %v", err)
	}
	_ = resp.Body.Close()

	// Wait for result
	select {
	case result := <-resultChan:
		if result.Username != "testuser" {
			t.Errorf("Username = %q, want %q", result.Username, "testuser")
		}
		if result.UserID != "12345" {
			t.Errorf("UserID = %q, want %q", result.UserID, "12345")
		}
		if result.Token != "long-token" {
			t.Errorf("Token = %q, want %q", result.Token, "long-token")
		}
		if result.ExpiresIn != 5184000 {
			t.Errorf("ExpiresIn = %d, want %d", result.ExpiresIn, 5184000)
		}
	case err := <-errChan:
		t.Fatalf("OAuth flow error: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for OAuth flow")
	}
}

func TestOAuthFlow_Run_NoCode(t *testing.T) {
	keychain := NewMockKeychain()
	mockClient := &api.MockClient{}

	flow := &OAuthFlow{
		AppID:       "test-app-id",
		AppSecret:   "test-app-secret",
		Client:      mockClient,
		Keychain:    keychain,
		OpenBrowser: noopBrowser,
	}

	errChan := make(chan error, 1)

	go func() {
		_, err := flow.Run()
		errChan <- err
	}()

	time.Sleep(200 * time.Millisecond)

	// Send callback without code
	resp, err := http.Get("http://localhost:8080/callback?error=access_denied")
	if err != nil {
		t.Fatalf("callback request error: %v", err)
	}
	_ = resp.Body.Close()

	select {
	case err := <-errChan:
		if err == nil {
			t.Error("expected error for missing code")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

func TestOAuthFlow_Run_ExchangeCodeFails(t *testing.T) {
	keychain := NewMockKeychain()
	mockClient := &api.MockClient{
		ExchangeCodeForTokenFn: func(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error) {
			return nil, fmt.Errorf("exchange failed")
		},
	}

	flow := &OAuthFlow{
		AppID:       "test-app-id",
		AppSecret:   "test-app-secret",
		Client:      mockClient,
		Keychain:    keychain,
		OpenBrowser: noopBrowser,
	}

	errChan := make(chan error, 1)

	go func() {
		_, err := flow.Run()
		errChan <- err
	}()

	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/callback?code=test-code")
	if err != nil {
		t.Fatalf("callback request error: %v", err)
	}
	_ = resp.Body.Close()

	select {
	case err := <-errChan:
		if err == nil {
			t.Error("expected error for failed exchange")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

func TestOAuthFlow_Run_LongLivedTokenFails(t *testing.T) {
	keychain := NewMockKeychain()
	mockClient := &api.MockClient{
		ExchangeCodeForTokenFn: func(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error) {
			return &models.TokenResponse{AccessToken: "short"}, nil
		},
		ExchangeForLongLivedTokenFn: func(appID, appSecret, shortToken string) (*models.TokenResponse, error) {
			return nil, fmt.Errorf("long-lived exchange failed")
		},
	}

	flow := &OAuthFlow{
		AppID:       "test-app-id",
		AppSecret:   "test-app-secret",
		Client:      mockClient,
		Keychain:    keychain,
		OpenBrowser: noopBrowser,
	}

	errChan := make(chan error, 1)

	go func() {
		_, err := flow.Run()
		errChan <- err
	}()

	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/callback?code=test-code")
	if err != nil {
		t.Fatalf("callback request error: %v", err)
	}
	_ = resp.Body.Close()

	select {
	case err := <-errChan:
		if err == nil {
			t.Error("expected error for failed long-lived token exchange")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

func TestOAuthFlow_Run_GetProfileFails(t *testing.T) {
	keychain := NewMockKeychain()
	mockClient := &api.MockClient{
		ExchangeCodeForTokenFn: func(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error) {
			return &models.TokenResponse{AccessToken: "short"}, nil
		},
		ExchangeForLongLivedTokenFn: func(appID, appSecret, shortToken string) (*models.TokenResponse, error) {
			return &models.TokenResponse{AccessToken: "long", ExpiresIn: 5184000}, nil
		},
		GetUserProfileFn: func(token string) (*models.UserProfile, error) {
			return nil, fmt.Errorf("profile fetch failed")
		},
	}

	flow := &OAuthFlow{
		AppID:       "test-app-id",
		AppSecret:   "test-app-secret",
		Client:      mockClient,
		Keychain:    keychain,
		OpenBrowser: noopBrowser,
	}

	errChan := make(chan error, 1)

	go func() {
		_, err := flow.Run()
		errChan <- err
	}()

	time.Sleep(200 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/callback?code=test-code")
	if err != nil {
		t.Fatalf("callback request error: %v", err)
	}
	_ = resp.Body.Close()

	select {
	case err := <-errChan:
		if err == nil {
			t.Error("expected error for failed profile fetch")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

func TestOAuthFlow_Run_NoCodeNoError(t *testing.T) {
	keychain := NewMockKeychain()
	mockClient := &api.MockClient{}

	flow := &OAuthFlow{
		AppID:       "test-app-id",
		AppSecret:   "test-app-secret",
		Client:      mockClient,
		Keychain:    keychain,
		OpenBrowser: noopBrowser,
	}

	errChan := make(chan error, 1)

	go func() {
		_, err := flow.Run()
		errChan <- err
	}()

	time.Sleep(200 * time.Millisecond)

	// Send callback without code or error
	resp, err := http.Get("http://localhost:8080/callback")
	if err != nil {
		t.Fatalf("callback request error: %v", err)
	}
	_ = resp.Body.Close()

	select {
	case err := <-errChan:
		if err == nil {
			t.Error("expected error for empty callback")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}
