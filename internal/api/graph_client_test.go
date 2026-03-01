package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pj4533/ig-cli/internal/models"
)

// newTestGraphClient creates a GraphClient pointed at a test server.
func newTestGraphClient(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *GraphClient) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := &GraphClient{
		httpClient: server.Client(),
		token:      "test-token",
		baseURL:    server.URL,
		fbBaseURL:  server.URL,
	}
	return server, client
}

func TestGraphClient_ListMedia_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/12345/media") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := models.PaginatedResponse[models.Media]{
			Data: []models.Media{
				{ID: "m1", Caption: "Test", MediaType: "IMAGE", LikeCount: 42},
				{ID: "m2", Caption: "Test 2", MediaType: "VIDEO", LikeCount: 10},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	media, err := client.ListMedia("12345", 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(media) != 2 {
		t.Fatalf("length = %d, want 2", len(media))
	}
	if media[0].ID != "m1" {
		t.Errorf("ID = %q, want %q", media[0].ID, "m1")
	}
}

func TestGraphClient_ListMedia_WithLimit(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		limit := r.URL.Query().Get("limit")
		if limit != "5" {
			t.Errorf("limit = %q, want %q", limit, "5")
		}
		resp := models.PaginatedResponse[models.Media]{
			Data: []models.Media{{ID: "m1"}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	media, err := client.ListMedia("12345", 5)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(media) != 1 {
		t.Errorf("length = %d, want 1", len(media))
	}
}

func TestGraphClient_GetMediaInsights_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/m1/insights") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := models.PaginatedResponse[models.MediaInsight]{
			Data: []models.MediaInsight{
				{Name: "impressions", Values: []models.MetricValue{{Value: 5000}}},
				{Name: "reach", Values: []models.MetricValue{{Value: 3000}}},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	insights, err := client.GetMediaInsights("m1")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(insights) != 2 {
		t.Fatalf("length = %d, want 2", len(insights))
	}
	if insights[0].Name != "impressions" {
		t.Errorf("Name = %q, want %q", insights[0].Name, "impressions")
	}
}

func TestGraphClient_ListComments_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/m1/comments") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := models.PaginatedResponse[models.Comment]{
			Data: []models.Comment{
				{ID: "c1", Text: "Great!", Username: "user1"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	comments, err := client.ListComments("m1", 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("length = %d, want 1", len(comments))
	}
	if comments[0].Username != "user1" {
		t.Errorf("Username = %q, want %q", comments[0].Username, "user1")
	}
}

func TestGraphClient_ListComments_WithLimit(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		resp := models.PaginatedResponse[models.Comment]{
			Data: []models.Comment{{ID: "c1"}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	comments, err := client.ListComments("m1", 10)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(comments) != 1 {
		t.Errorf("length = %d, want 1", len(comments))
	}
}

func TestGraphClient_ListReplies_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/c1/replies") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		resp := models.PaginatedResponse[models.Comment]{
			Data: []models.Comment{
				{ID: "r1", Text: "Thanks!", Username: "author"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	replies, err := client.ListReplies("c1", 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(replies) != 1 {
		t.Fatalf("length = %d, want 1", len(replies))
	}
}

func TestGraphClient_ListReplies_WithLimit(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		resp := models.PaginatedResponse[models.Comment]{
			Data: []models.Comment{{ID: "r1"}},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	replies, err := client.ListReplies("c1", 5)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(replies) != 1 {
		t.Errorf("length = %d, want 1", len(replies))
	}
}

func TestGraphClient_GetAccountInsights_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if period != "week" {
			t.Errorf("period = %q, want %q", period, "week")
		}
		resp := models.PaginatedResponse[models.AccountInsight]{
			Data: []models.AccountInsight{
				{Name: "impressions", Period: "week", Values: []models.MetricValue{{Value: 10000}}},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	insights, err := client.GetAccountInsights("12345", "week")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(insights) != 1 {
		t.Fatalf("length = %d, want 1", len(insights))
	}
}

func TestGraphClient_GetAccountInsights_DefaultPeriod(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if period != "day" {
			t.Errorf("period = %q, want %q (default)", period, "day")
		}
		resp := models.PaginatedResponse[models.AccountInsight]{
			Data: []models.AccountInsight{},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	_, err := client.GetAccountInsights("12345", "")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}

func TestGraphClient_GetAudienceDemographics_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		period := r.URL.Query().Get("period")
		if period != "lifetime" {
			t.Errorf("period = %q, want %q", period, "lifetime")
		}
		resp := models.PaginatedResponse[models.AudienceDemographic]{
			Data: []models.AudienceDemographic{
				{Name: "audience_city", Period: "lifetime"},
				{Name: "audience_country", Period: "lifetime"},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	demographics, err := client.GetAudienceDemographics("12345")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(demographics) != 2 {
		t.Fatalf("length = %d, want 2", len(demographics))
	}
}

func TestGraphClient_DiscoverUser_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		fields := r.URL.Query().Get("fields")
		if !strings.Contains(fields, "business_discovery") {
			t.Errorf("fields should contain business_discovery: %q", fields)
		}
		if !strings.Contains(fields, "targetuser") {
			t.Errorf("fields should contain username: %q", fields)
		}
		result := struct {
			BusinessDiscovery models.BusinessDiscovery `json:"business_discovery"`
		}{
			BusinessDiscovery: models.BusinessDiscovery{
				ID:             "bd1",
				Username:       "targetuser",
				FollowersCount: 5000,
			},
		}
		_ = json.NewEncoder(w).Encode(result)
	})
	defer server.Close()

	discovery, err := client.DiscoverUser("12345", "targetuser")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if discovery.Username != "targetuser" {
		t.Errorf("Username = %q, want %q", discovery.Username, "targetuser")
	}
	if discovery.FollowersCount != 5000 {
		t.Errorf("FollowersCount = %d, want %d", discovery.FollowersCount, 5000)
	}
}

func TestGraphClient_ExchangeCodeForToken_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("method = %q, want POST", r.Method)
		}
		resp := models.TokenResponse{
			AccessToken: "short-lived-token",
			TokenType:   "bearer",
			ExpiresIn:   3600,
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	resp, err := client.ExchangeCodeForToken("app-id", "secret", "http://localhost/cb", "code123")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.AccessToken != "short-lived-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "short-lived-token")
	}
}

func TestGraphClient_ExchangeForLongLivedToken_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		grantType := r.URL.Query().Get("grant_type")
		if grantType != "ig_exchange_token" {
			t.Errorf("grant_type = %q, want %q", grantType, "ig_exchange_token")
		}
		resp := models.TokenResponse{
			AccessToken: "long-lived-token",
			ExpiresIn:   5184000,
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	resp, err := client.ExchangeForLongLivedToken("app-id", "secret", "short-token")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.AccessToken != "long-lived-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "long-lived-token")
	}
}

func TestGraphClient_RefreshLongLivedToken_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		grantType := r.URL.Query().Get("grant_type")
		if grantType != "ig_refresh_token" {
			t.Errorf("grant_type = %q, want %q", grantType, "ig_refresh_token")
		}
		resp := models.TokenResponse{
			AccessToken: "refreshed-token",
			ExpiresIn:   5184000,
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	resp, err := client.RefreshLongLivedToken("old-token")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if resp.AccessToken != "refreshed-token" {
		t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "refreshed-token")
	}
}

func TestGraphClient_GetUserProfile_Real(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/me") {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		profile := models.UserProfile{
			ID:       "12345",
			Username: "testuser",
			Name:     "Test User",
		}
		_ = json.NewEncoder(w).Encode(profile)
	})
	defer server.Close()

	profile, err := client.GetUserProfile("test-token")
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if profile.Username != "testuser" {
		t.Errorf("Username = %q, want %q", profile.Username, "testuser")
	}
}

func TestGraphClient_ErrorParsing(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		resp := models.APIErrorResponse{}
		resp.Error.Message = "Invalid OAuth access token"
		resp.Error.Type = "OAuthException"
		resp.Error.Code = 190
		resp.Error.FBTraceID = "trace123"
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	_, _, err := client.doRequest("GET", server.URL+"/test")
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if !apiErr.IsAuthExpired() {
		t.Error("expected IsAuthExpired() to be true for code 190")
	}
	if apiErr.FBTraceID != "trace123" {
		t.Errorf("FBTraceID = %q, want %q", apiErr.FBTraceID, "trace123")
	}
}

func TestGraphClient_RateLimitError(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		resp := models.APIErrorResponse{}
		resp.Error.Message = "Rate limit exceeded"
		resp.Error.Type = "OAuthException"
		resp.Error.Code = 4
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	_, _, err := client.doRequest("GET", server.URL+"/test")
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}

	if !apiErr.IsRateLimited() {
		t.Error("expected IsRateLimited() to be true for code 4")
	}
}

func TestBuildFacebookURL(t *testing.T) {
	client := NewGraphClient("test-token")
	url := client.buildFacebookURL("/oauth/access_token", nil)
	if !strings.Contains(url, FacebookBaseURL) {
		t.Errorf("URL should contain FacebookBaseURL: %q", url)
	}
	if !strings.Contains(url, "/oauth/access_token") {
		t.Errorf("URL should contain path: %q", url)
	}
}

func TestBuildFacebookURL_EmptyFBBase(t *testing.T) {
	client := &GraphClient{fbBaseURL: ""}
	url := client.buildFacebookURL("/test", nil)
	if !strings.Contains(url, FacebookBaseURL) {
		t.Errorf("URL should fall back to FacebookBaseURL: %q", url)
	}
}

func TestBuildURL_WithParams(t *testing.T) {
	client := NewGraphClient("my-token")
	url := client.buildURL("/12345/media", nil)

	if !strings.Contains(url, "access_token=my-token") {
		t.Errorf("URL should contain access_token: %q", url)
	}
	if !strings.Contains(url, BaseURL) {
		t.Errorf("URL should contain BaseURL: %q", url)
	}
	if !strings.Contains(url, "/12345/media") {
		t.Errorf("URL should contain path: %q", url)
	}
}

func TestBuildURL_EmptyToken(t *testing.T) {
	client := NewGraphClient("")
	url := client.buildURL("/test", nil)

	if strings.Contains(url, "access_token") {
		t.Errorf("URL should not contain access_token when token is empty: %q", url)
	}
}

func TestGraphClient_APIError_Methods(t *testing.T) {
	// Test various error scenarios on real endpoint methods
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(`{"error":{"message":"Server error","type":"ServerException","code":2}}`))
	})
	defer server.Close()

	_, err := client.ListMedia("12345", 0)
	if err == nil {
		t.Error("ListMedia: expected error")
	}

	_, err = client.GetMediaInsights("m1")
	if err == nil {
		t.Error("GetMediaInsights: expected error")
	}

	_, err = client.ListComments("m1", 0)
	if err == nil {
		t.Error("ListComments: expected error")
	}

	_, err = client.ListReplies("c1", 0)
	if err == nil {
		t.Error("ListReplies: expected error")
	}

	_, err = client.GetAccountInsights("12345", "day")
	if err == nil {
		t.Error("GetAccountInsights: expected error")
	}

	_, err = client.GetAudienceDemographics("12345")
	if err == nil {
		t.Error("GetAudienceDemographics: expected error")
	}

	_, err = client.DiscoverUser("12345", "target")
	if err == nil {
		t.Error("DiscoverUser: expected error")
	}

	_, err = client.ExchangeCodeForToken("a", "b", "c", "d")
	if err == nil {
		t.Error("ExchangeCodeForToken: expected error")
	}

	_, err = client.ExchangeForLongLivedToken("a", "b", "c")
	if err == nil {
		t.Error("ExchangeForLongLivedToken: expected error")
	}

	_, err = client.RefreshLongLivedToken("t")
	if err == nil {
		t.Error("RefreshLongLivedToken: expected error")
	}

	_, err = client.GetUserProfile("t")
	if err == nil {
		t.Error("GetUserProfile: expected error")
	}
}

func TestGraphClient_InvalidJSON_Responses(t *testing.T) {
	server, client := newTestGraphClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not valid json`))
	})
	defer server.Close()

	_, err := client.GetMediaInsights("m1")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	_, err = client.GetAccountInsights("12345", "day")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	_, err = client.GetAudienceDemographics("12345")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	_, err = client.DiscoverUser("12345", "target")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	_, err = client.ExchangeCodeForToken("a", "b", "c", "d")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	_, err = client.ExchangeForLongLivedToken("a", "b", "c")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	_, err = client.RefreshLongLivedToken("t")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	_, err = client.GetUserProfile("t")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
