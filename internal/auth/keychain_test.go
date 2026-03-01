package auth

import (
	"testing"
)

func TestMockKeychain(t *testing.T) {
	kc := NewMockKeychain()

	// Set and Get
	if err := kc.Set("key1", "value1"); err != nil {
		t.Fatalf("Set error: %v", err)
	}

	val, err := kc.Get("key1")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if val != "value1" {
		t.Errorf("Get = %q, want %q", val, "value1")
	}

	// Get non-existent
	_, err = kc.Get("nonexistent")
	if err == nil {
		t.Error("Get should error for non-existent key")
	}

	// Delete
	if err := kc.Delete("key1"); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	_, err = kc.Get("key1")
	if err == nil {
		t.Error("Get should error after Delete")
	}

	// Delete non-existent (should not error)
	if err := kc.Delete("nonexistent"); err != nil {
		t.Errorf("Delete non-existent error: %v", err)
	}
}

func TestMockKeychainOverwrite(t *testing.T) {
	kc := NewMockKeychain()

	_ = kc.Set("key", "v1")
	_ = kc.Set("key", "v2")

	val, err := kc.Get("key")
	if err != nil {
		t.Fatalf("Get error: %v", err)
	}
	if val != "v2" {
		t.Errorf("Get = %q, want %q (should overwrite)", val, "v2")
	}
}

func TestTokenKey(t *testing.T) {
	if got := TokenKey("testuser"); got != "token:testuser" {
		t.Errorf("TokenKey = %q, want %q", got, "token:testuser")
	}
}

func TestUserIDKey(t *testing.T) {
	if got := UserIDKey("testuser"); got != "user_id:testuser" {
		t.Errorf("UserIDKey = %q, want %q", got, "user_id:testuser")
	}
}

func TestTokenExpiryKey(t *testing.T) {
	if got := TokenExpiryKey("testuser"); got != "token_expiry:testuser" {
		t.Errorf("TokenExpiryKey = %q, want %q", got, "token_expiry:testuser")
	}
}

func TestKeychainStoreInterface(t *testing.T) {
	// Verify MockKeychain satisfies KeychainStore
	var _ KeychainStore = &MockKeychain{}
	var _ KeychainStore = &OSKeychain{}
}
