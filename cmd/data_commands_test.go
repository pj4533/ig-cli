package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/pj4533/ig-cli/internal/api"
	"github.com/pj4533/ig-cli/internal/auth"
	"github.com/pj4533/ig-cli/internal/config"
	"github.com/pj4533/ig-cli/internal/models"
	"github.com/spf13/viper"
)

// setupTestEnv sets up a test environment with a mock client, keychain, and config.
func setupTestEnv(t *testing.T, mock *api.MockClient) func() {
	t.Helper()

	// Create temp home dir for config
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	// Reset viper
	viper.Reset()

	// Save test config
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

	// Reset viper for fresh load
	viper.Reset()

	// Override client factory
	origClientFactory := clientFactory
	clientFactory = func(token string) api.Client {
		return mock
	}

	// Override keychain factory with mock keychain that has a valid token
	origKeychainFactory := keychainFactory
	mockKeychain := auth.NewMockKeychain()
	_ = mockKeychain.Set(auth.TokenKey("testuser"), "test-token")
	expiry := time.Now().Add(30 * 24 * time.Hour).Unix()
	_ = mockKeychain.Set(auth.TokenExpiryKey("testuser"), strconv.FormatInt(expiry, 10))

	keychainFactory = func() auth.KeychainStore {
		return mockKeychain
	}

	// Reset account flag
	accountFlag = ""

	return func() {
		clientFactory = origClientFactory
		keychainFactory = origKeychainFactory
		viper.Reset()
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestRunMediaList(t *testing.T) {
	mock := &api.MockClient{
		ListMediaFn: func(userID string, limit int) ([]models.Media, error) {
			if userID != "12345" {
				t.Errorf("userID = %q, want %q", userID, "12345")
			}
			return []models.Media{
				{ID: "m1", Caption: "Test post", MediaType: "IMAGE", LikeCount: 42},
			}, nil
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	output := captureStdout(t, func() {
		err := runMediaList(nil, nil)
		if err != nil {
			t.Fatalf("runMediaList error: %v", err)
		}
	})

	var result []models.Media
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v\noutput: %s", err, output)
	}
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
	if result[0].ID != "m1" {
		t.Errorf("ID = %q, want %q", result[0].ID, "m1")
	}
}

func TestRunMediaList_Error(t *testing.T) {
	mock := &api.MockClient{
		ListMediaFn: func(userID string, limit int) ([]models.Media, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	err := runMediaList(nil, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRunMediaInsights(t *testing.T) {
	mock := &api.MockClient{
		GetMediaInsightsFn: func(mediaID string) ([]models.MediaInsight, error) {
			if mediaID != "m1" {
				t.Errorf("mediaID = %q, want %q", mediaID, "m1")
			}
			return []models.MediaInsight{
				{Name: "impressions", Values: []models.MetricValue{{Value: 1000}}},
			}, nil
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	output := captureStdout(t, func() {
		err := runMediaInsights(nil, []string{"m1"})
		if err != nil {
			t.Fatalf("runMediaInsights error: %v", err)
		}
	})

	var result []models.MediaInsight
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v\noutput: %s", err, output)
	}
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
}

func TestRunMediaInsights_Error(t *testing.T) {
	mock := &api.MockClient{
		GetMediaInsightsFn: func(mediaID string) ([]models.MediaInsight, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	err := runMediaInsights(nil, []string{"m1"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRunCommentsList(t *testing.T) {
	mock := &api.MockClient{
		ListCommentsFn: func(mediaID string, limit int) ([]models.Comment, error) {
			return []models.Comment{
				{ID: "c1", Text: "Great!", Username: "user1"},
			}, nil
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	output := captureStdout(t, func() {
		err := runCommentsList(nil, []string{"m1"})
		if err != nil {
			t.Fatalf("runCommentsList error: %v", err)
		}
	})

	var result []models.Comment
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
}

func TestRunCommentsList_Error(t *testing.T) {
	mock := &api.MockClient{
		ListCommentsFn: func(mediaID string, limit int) ([]models.Comment, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	err := runCommentsList(nil, []string{"m1"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRunCommentsReplies(t *testing.T) {
	mock := &api.MockClient{
		ListRepliesFn: func(commentID string, limit int) ([]models.Comment, error) {
			return []models.Comment{
				{ID: "r1", Text: "Thanks!", Username: "author"},
			}, nil
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	output := captureStdout(t, func() {
		err := runCommentsReplies(nil, []string{"c1"})
		if err != nil {
			t.Fatalf("runCommentsReplies error: %v", err)
		}
	})

	var result []models.Comment
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
}

func TestRunCommentsReplies_Error(t *testing.T) {
	mock := &api.MockClient{
		ListRepliesFn: func(commentID string, limit int) ([]models.Comment, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	err := runCommentsReplies(nil, []string{"c1"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestRunInsightsAccount(t *testing.T) {
	mock := &api.MockClient{
		GetAccountInsightsFn: func(userID string, period string) ([]models.AccountInsight, error) {
			return []models.AccountInsight{
				{Name: "impressions", Period: "day", Values: []models.MetricValue{{Value: 5000}}},
			}, nil
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	output := captureStdout(t, func() {
		err := runInsightsAccount(nil, nil)
		if err != nil {
			t.Fatalf("runInsightsAccount error: %v", err)
		}
	})

	var result []models.AccountInsight
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
}

func TestRunInsightsAccount_Error(t *testing.T) {
	mock := &api.MockClient{
		GetAccountInsightsFn: func(userID string, period string) ([]models.AccountInsight, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	err := runInsightsAccount(nil, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRunInsightsAudience(t *testing.T) {
	mock := &api.MockClient{
		GetAudienceDemographicsFn: func(userID string) ([]models.AudienceDemographic, error) {
			return []models.AudienceDemographic{
				{Name: "audience_city", Period: "lifetime"},
			}, nil
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	output := captureStdout(t, func() {
		err := runInsightsAudience(nil, nil)
		if err != nil {
			t.Fatalf("runInsightsAudience error: %v", err)
		}
	})

	var result []models.AudienceDemographic
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("result length = %d, want 1", len(result))
	}
}

func TestRunInsightsAudience_Error(t *testing.T) {
	mock := &api.MockClient{
		GetAudienceDemographicsFn: func(userID string) ([]models.AudienceDemographic, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	err := runInsightsAudience(nil, nil)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRunDiscover(t *testing.T) {
	mock := &api.MockClient{
		DiscoverUserFn: func(userID string, username string) (*models.BusinessDiscovery, error) {
			if username != "targetuser" {
				t.Errorf("username = %q, want %q", username, "targetuser")
			}
			return &models.BusinessDiscovery{
				ID:             "bd1",
				Username:       "targetuser",
				FollowersCount: 5000,
			}, nil
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	output := captureStdout(t, func() {
		err := runDiscover(nil, []string{"targetuser"})
		if err != nil {
			t.Fatalf("runDiscover error: %v", err)
		}
	})

	var result models.BusinessDiscovery
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if result.Username != "targetuser" {
		t.Errorf("Username = %q, want %q", result.Username, "targetuser")
	}
}

func TestRunDiscover_Error(t *testing.T) {
	mock := &api.MockClient{
		DiscoverUserFn: func(userID string, username string) (*models.BusinessDiscovery, error) {
			return nil, fmt.Errorf("API error")
		},
	}

	cleanup := setupTestEnv(t, mock)
	defer cleanup()

	err := runDiscover(nil, []string{"targetuser"})
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetClient_NoConfig(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("HOME", dir)

	viper.Reset()
	defer viper.Reset()

	_, _, err := getClient()
	if err == nil {
		t.Error("expected error when no config")
	}
}
