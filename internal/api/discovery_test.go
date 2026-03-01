package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pj4533/ig-cli/internal/models"
)

func TestDiscoverUser_Mock(t *testing.T) {
	mock := &MockClient{
		DiscoverUserFn: func(userID string, username string) (*models.BusinessDiscovery, error) {
			if userID != "12345" {
				t.Errorf("userID = %q, want %q", userID, "12345")
			}
			if username != "targetuser" {
				t.Errorf("username = %q, want %q", username, "targetuser")
			}
			return &models.BusinessDiscovery{
				ID:             "bd1",
				Username:       "targetuser",
				Name:           "Target User",
				FollowersCount: 10000,
				MediaCount:     500,
			}, nil
		},
	}

	discovery, err := mock.DiscoverUser("12345", "targetuser")
	if err != nil {
		t.Fatalf("DiscoverUser error: %v", err)
	}
	if discovery.Username != "targetuser" {
		t.Errorf("Username = %q, want %q", discovery.Username, "targetuser")
	}
	if discovery.FollowersCount != 10000 {
		t.Errorf("FollowersCount = %d, want %d", discovery.FollowersCount, 10000)
	}
}

func TestDiscoverUser_HTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		result := struct {
			BusinessDiscovery models.BusinessDiscovery `json:"business_discovery"`
		}{
			BusinessDiscovery: models.BusinessDiscovery{
				ID:             "bd1",
				Username:       "targetuser",
				Name:           "Target",
				FollowersCount: 5000,
			},
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(result)
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

	var result struct {
		BusinessDiscovery models.BusinessDiscovery `json:"business_discovery"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if result.BusinessDiscovery.Username != "targetuser" {
		t.Errorf("Username = %q, want %q", result.BusinessDiscovery.Username, "targetuser")
	}
}
