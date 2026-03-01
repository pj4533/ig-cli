package auth

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/pj4533/ig-cli/internal/api"
	"github.com/pj4533/ig-cli/internal/models"
)

func TestGetValidToken_Fresh(t *testing.T) {
	kc := NewMockKeychain()
	_ = kc.Set(TokenKey("testuser"), "my-token")
	// Set expiry far in the future
	expiry := time.Now().Add(30 * 24 * time.Hour).Unix()
	_ = kc.Set(TokenExpiryKey("testuser"), strconv.FormatInt(expiry, 10))

	client := &api.MockClient{}
	tm := NewTokenManager(kc, client)

	token, err := tm.GetValidToken("testuser")
	if err != nil {
		t.Fatalf("GetValidToken error: %v", err)
	}
	if token != "my-token" {
		t.Errorf("token = %q, want %q", token, "my-token")
	}
}

func TestGetValidToken_NearExpiry_Refreshes(t *testing.T) {
	kc := NewMockKeychain()
	_ = kc.Set(TokenKey("testuser"), "old-token")
	// Set expiry within 24 hours
	expiry := time.Now().Add(12 * time.Hour).Unix()
	_ = kc.Set(TokenExpiryKey("testuser"), strconv.FormatInt(expiry, 10))

	refreshCalled := false
	client := &api.MockClient{
		RefreshLongLivedTokenFn: func(token string) (*models.TokenResponse, error) {
			refreshCalled = true
			if token != "old-token" {
				t.Errorf("refresh called with %q, want %q", token, "old-token")
			}
			return &models.TokenResponse{
				AccessToken: "new-token",
				ExpiresIn:   5184000,
			}, nil
		},
	}

	tm := NewTokenManager(kc, client)

	token, err := tm.GetValidToken("testuser")
	if err != nil {
		t.Fatalf("GetValidToken error: %v", err)
	}
	if !refreshCalled {
		t.Error("refresh was not called")
	}
	if token != "new-token" {
		t.Errorf("token = %q, want %q", token, "new-token")
	}

	// Verify new token is stored
	stored, _ := kc.Get(TokenKey("testuser"))
	if stored != "new-token" {
		t.Errorf("stored token = %q, want %q", stored, "new-token")
	}
}

func TestGetValidToken_RefreshFails_ReturnsOld(t *testing.T) {
	kc := NewMockKeychain()
	_ = kc.Set(TokenKey("testuser"), "old-token")
	expiry := time.Now().Add(12 * time.Hour).Unix()
	_ = kc.Set(TokenExpiryKey("testuser"), strconv.FormatInt(expiry, 10))

	client := &api.MockClient{
		RefreshLongLivedTokenFn: func(token string) (*models.TokenResponse, error) {
			return nil, fmt.Errorf("refresh failed")
		},
	}

	tm := NewTokenManager(kc, client)

	token, err := tm.GetValidToken("testuser")
	if err != nil {
		t.Fatalf("GetValidToken error: %v", err)
	}
	if token != "old-token" {
		t.Errorf("token = %q, want %q (should fall back to old)", token, "old-token")
	}
}

func TestGetValidToken_NoExpiry(t *testing.T) {
	kc := NewMockKeychain()
	_ = kc.Set(TokenKey("testuser"), "my-token")
	// No expiry set

	client := &api.MockClient{}
	tm := NewTokenManager(kc, client)

	token, err := tm.GetValidToken("testuser")
	if err != nil {
		t.Fatalf("GetValidToken error: %v", err)
	}
	if token != "my-token" {
		t.Errorf("token = %q, want %q", token, "my-token")
	}
}

func TestGetValidToken_NoToken(t *testing.T) {
	kc := NewMockKeychain()
	client := &api.MockClient{}
	tm := NewTokenManager(kc, client)

	_, err := tm.GetValidToken("testuser")
	if err == nil {
		t.Error("expected error for missing token")
	}
}

func TestStoreToken(t *testing.T) {
	kc := NewMockKeychain()
	client := &api.MockClient{}
	tm := NewTokenManager(kc, client)

	if err := tm.StoreToken("testuser", "my-token", 5184000); err != nil {
		t.Fatalf("StoreToken error: %v", err)
	}

	token, err := kc.Get(TokenKey("testuser"))
	if err != nil {
		t.Fatalf("Get token error: %v", err)
	}
	if token != "my-token" {
		t.Errorf("token = %q, want %q", token, "my-token")
	}

	expiryStr, err := kc.Get(TokenExpiryKey("testuser"))
	if err != nil {
		t.Fatalf("Get expiry error: %v", err)
	}

	expiry, err := strconv.ParseInt(expiryStr, 10, 64)
	if err != nil {
		t.Fatalf("parse expiry error: %v", err)
	}

	// Expiry should be roughly now + 5184000 seconds
	expected := time.Now().Add(5184000 * time.Second).Unix()
	if diff := expected - expiry; diff < -5 || diff > 5 {
		t.Errorf("expiry diff = %d seconds, want within 5 seconds", diff)
	}
}

func TestGetValidToken_InvalidExpiryFormat(t *testing.T) {
	kc := NewMockKeychain()
	_ = kc.Set(TokenKey("testuser"), "my-token")
	_ = kc.Set(TokenExpiryKey("testuser"), "not-a-number")

	client := &api.MockClient{}
	tm := NewTokenManager(kc, client)

	token, err := tm.GetValidToken("testuser")
	if err != nil {
		t.Fatalf("GetValidToken error: %v", err)
	}
	if token != "my-token" {
		t.Errorf("token = %q, want %q (should use as-is with invalid expiry)", token, "my-token")
	}
}
