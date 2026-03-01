package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pj4533/ig-cli/internal/models"
)

func TestExchangeCodeForToken_Mock(t *testing.T) {
	mock := &MockClient{
		ExchangeCodeForTokenFn: func(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error) {
			if code != "test-code" {
				t.Errorf("code = %q, want %q", code, "test-code")
			}
			return &models.TokenResponse{
				AccessToken: "short-lived-token",
				TokenType:   "bearer",
				ExpiresIn:   3600,
			}, nil
		},
	}

	resp, err := mock.ExchangeCodeForToken("app-id", "app-secret", "http://localhost:8080/callback", "test-code")
	if err != nil {
		t.Fatalf("ExchangeCodeForToken error: %v", err)
	}
	if resp.AccessToken != "short-lived-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "short-lived-token")
	}
}

func TestExchangeForLongLivedToken_Mock(t *testing.T) {
	mock := &MockClient{
		ExchangeForLongLivedTokenFn: func(appID, appSecret, shortToken string) (*models.TokenResponse, error) {
			return &models.TokenResponse{
				AccessToken: "long-lived-token",
				ExpiresIn:   5184000,
			}, nil
		},
	}

	resp, err := mock.ExchangeForLongLivedToken("app-id", "app-secret", "short-token")
	if err != nil {
		t.Fatalf("ExchangeForLongLivedToken error: %v", err)
	}
	if resp.AccessToken != "long-lived-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "long-lived-token")
	}
	if resp.ExpiresIn != 5184000 {
		t.Errorf("ExpiresIn = %d, want %d", resp.ExpiresIn, 5184000)
	}
}

func TestRefreshLongLivedToken_Mock(t *testing.T) {
	mock := &MockClient{
		RefreshLongLivedTokenFn: func(token string) (*models.TokenResponse, error) {
			return &models.TokenResponse{
				AccessToken: "refreshed-token",
				ExpiresIn:   5184000,
			}, nil
		},
	}

	resp, err := mock.RefreshLongLivedToken("old-token")
	if err != nil {
		t.Fatalf("RefreshLongLivedToken error: %v", err)
	}
	if resp.AccessToken != "refreshed-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "refreshed-token")
	}
}

func TestGetUserProfile_Mock(t *testing.T) {
	mock := &MockClient{
		GetUserProfileFn: func(token string) (*models.UserProfile, error) {
			return &models.UserProfile{
				ID:       "12345",
				Username: "testuser",
				Name:     "Test User",
			}, nil
		},
	}

	profile, err := mock.GetUserProfile("token")
	if err != nil {
		t.Fatalf("GetUserProfile error: %v", err)
	}
	if profile.Username != "testuser" {
		t.Errorf("Username = %q, want %q", profile.Username, "testuser")
	}
}

func TestGetUserProfile_HTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		profile := models.UserProfile{
			ID:       "12345",
			Username: "testuser",
			Name:     "Test User",
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(profile)
	}))
	defer server.Close()

	client := &GraphClient{
		httpClient: server.Client(),
		token:      "test-token",
	}

	body, _, err := client.doRequest("GET", server.URL)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var profile models.UserProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if profile.Username != "testuser" {
		t.Errorf("Username = %q, want %q", profile.Username, "testuser")
	}
}

func TestExchangeCodeForToken_HTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.TokenResponse{
			AccessToken: "test-token",
			TokenType:   "bearer",
			ExpiresIn:   3600,
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &GraphClient{
		httpClient: server.Client(),
		token:      "",
	}

	body, _, err := client.doRequest("POST", server.URL)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	var resp models.TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if resp.AccessToken != "test-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "test-token")
	}
}

func TestMockClient_NilFunctions(t *testing.T) {
	mock := &MockClient{}

	// All nil functions should return nil, nil
	media, err := mock.ListMedia("id", 0)
	if err != nil || media != nil {
		t.Errorf("ListMedia = (%v, %v), want (nil, nil)", media, err)
	}

	insights, err := mock.GetMediaInsights("id")
	if err != nil || insights != nil {
		t.Errorf("GetMediaInsights = (%v, %v), want (nil, nil)", insights, err)
	}

	comments, err := mock.ListComments("id", 0)
	if err != nil || comments != nil {
		t.Errorf("ListComments = (%v, %v), want (nil, nil)", comments, err)
	}

	replies, err := mock.ListReplies("id", 0)
	if err != nil || replies != nil {
		t.Errorf("ListReplies = (%v, %v), want (nil, nil)", replies, err)
	}

	acctInsights, err := mock.GetAccountInsights("id", "day")
	if err != nil || acctInsights != nil {
		t.Errorf("GetAccountInsights = (%v, %v), want (nil, nil)", acctInsights, err)
	}

	demographics, err := mock.GetAudienceDemographics("id")
	if err != nil || demographics != nil {
		t.Errorf("GetAudienceDemographics = (%v, %v), want (nil, nil)", demographics, err)
	}

	discovery, err := mock.DiscoverUser("id", "user")
	if err != nil || discovery != nil {
		t.Errorf("DiscoverUser = (%v, %v), want (nil, nil)", discovery, err)
	}

	token, err := mock.ExchangeCodeForToken("a", "b", "c", "d")
	if err != nil || token != nil {
		t.Errorf("ExchangeCodeForToken = (%v, %v), want (nil, nil)", token, err)
	}

	longToken, err := mock.ExchangeForLongLivedToken("a", "b", "c")
	if err != nil || longToken != nil {
		t.Errorf("ExchangeForLongLivedToken = (%v, %v), want (nil, nil)", longToken, err)
	}

	refreshed, err := mock.RefreshLongLivedToken("t")
	if err != nil || refreshed != nil {
		t.Errorf("RefreshLongLivedToken = (%v, %v), want (nil, nil)", refreshed, err)
	}

	profile, err := mock.GetUserProfile("t")
	if err != nil || profile != nil {
		t.Errorf("GetUserProfile = (%v, %v), want (nil, nil)", profile, err)
	}
}
