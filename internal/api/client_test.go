package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pj4533/ig-cli/internal/models"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *GraphClient) {
	t.Helper()
	server := httptest.NewServer(handler)
	client := &GraphClient{
		httpClient: server.Client(),
		token:      "test-token",
	}
	return server, client
}

func TestDoRequest_Success(t *testing.T) {
	server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"ok": true}`))
	})
	defer server.Close()

	body, _, err := client.doRequest("GET", server.URL+"/test")
	if err != nil {
		t.Fatalf("doRequest error: %v", err)
	}

	var result map[string]bool
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}
	if !result["ok"] {
		t.Error("expected ok=true")
	}
}

func TestDoRequest_APIError(t *testing.T) {
	server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		_, _ = w.Write([]byte(`{"error":{"message":"Bad request","type":"OAuthException","code":190,"fbtrace_id":"abc"}}`))
	})
	defer server.Close()

	_, _, err := client.doRequest("GET", server.URL+"/test")
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Code != 190 {
		t.Errorf("Code = %d, want 190", apiErr.Code)
	}
	if apiErr.Message != "Bad request" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "Bad request")
	}
}

func TestDoRequest_NonJSONError(t *testing.T) {
	server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Internal Server Error"))
	})
	defer server.Close()

	_, _, err := client.doRequest("GET", server.URL+"/test")
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("StatusCode = %d, want 500", apiErr.StatusCode)
	}
}

func TestDoRequest_RateLimitHeaders(t *testing.T) {
	server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-App-Usage", `{"call_count":85,"total_cputime":2,"total_time":5}`)
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{}`))
	})
	defer server.Close()

	_, rateLimit, err := client.doRequest("GET", server.URL+"/test")
	if err != nil {
		t.Fatalf("doRequest error: %v", err)
	}
	if rateLimit == nil {
		t.Fatal("expected rate limit info")
	}
	if rateLimit.Usage != 85 {
		t.Errorf("Usage = %d, want 85", rateLimit.Usage)
	}
}

func TestAutoPaginate(t *testing.T) {
	callCount := 0
	var serverURL string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var resp models.PaginatedResponse[models.Media]
		if callCount == 1 {
			resp.Data = []models.Media{{ID: "media-1"}}
			resp.Paging.Next = serverURL + "/next?page=2&access_token=test-token"
		} else {
			resp.Data = []models.Media{{ID: "media-2"}}
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()
	serverURL = server.URL

	client := &GraphClient{
		httpClient: server.Client(),
		token:      "test-token",
	}

	results, err := autoPaginate[models.Media](client, server.URL+"/start?page=1&access_token=test-token", 0)
	if err != nil {
		t.Fatalf("autoPaginate error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("results length = %d, want 2", len(results))
	}
}

func TestAutoPaginate_WithLimit(t *testing.T) {
	server, client := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		resp := models.PaginatedResponse[models.Media]{
			Data: []models.Media{{ID: "1"}, {ID: "2"}, {ID: "3"}},
		}
		resp.Paging.Next = r.URL.Scheme + "://" + r.Host + "/next"
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer server.Close()

	results, err := autoPaginate[models.Media](client, server.URL+"/start", 2)
	if err != nil {
		t.Fatalf("autoPaginate error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("results length = %d, want 2 (limit)", len(results))
	}
}

func TestParseRateLimitHeaders(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   *RateLimitInfo
	}{
		{
			name:   "valid",
			header: `{"call_count":50,"total_cputime":2,"total_time":5}`,
			want:   &RateLimitInfo{Usage: 50, Limit: 100},
		},
		{
			name:   "empty",
			header: "",
			want:   nil,
		},
		{
			name:   "invalid json",
			header: "not-json",
			want:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := http.Header{}
			if tt.header != "" {
				h.Set("X-App-Usage", tt.header)
			}
			got := parseRateLimitHeaders(h)
			if tt.want == nil {
				if got != nil {
					t.Errorf("got %+v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("got nil, want non-nil")
			}
			if got.Usage != tt.want.Usage {
				t.Errorf("Usage = %d, want %d", got.Usage, tt.want.Usage)
			}
		})
	}
}

func TestBuildURL(t *testing.T) {
	client := NewGraphClient("my-token")
	url := client.buildURL("/12345/media", nil)

	if url == "" {
		t.Fatal("buildURL returned empty string")
	}
	// Should contain the base URL, path, and token
	if len(url) < 10 {
		t.Errorf("URL seems too short: %q", url)
	}
}

func TestNewGraphClient(t *testing.T) {
	client := NewGraphClient("test-token")
	if client == nil {
		t.Fatal("NewGraphClient returned nil")
	}
	if client.token != "test-token" {
		t.Errorf("token = %q, want %q", client.token, "test-token")
	}
	if client.httpClient == nil {
		t.Error("httpClient is nil")
	}
}

func TestMockClient_Implements_Interface(t *testing.T) {
	var _ Client = &MockClient{}
	var _ Client = &GraphClient{}
}
