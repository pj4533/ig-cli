package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/pj4533/ig-cli/internal/models"
)

const (
	// BaseURL is the Instagram Graph API base URL.
	BaseURL = "https://graph.instagram.com/v22.0"

	// FacebookBaseURL is the Facebook Graph API base URL (used for token exchange).
	FacebookBaseURL = "https://graph.facebook.com/v22.0"
)

// Client defines the interface for Instagram Graph API operations.
type Client interface {
	// Media
	ListMedia(userID string, limit int) ([]models.Media, error)
	GetMediaInsights(mediaID string) ([]models.MediaInsight, error)

	// Comments
	ListComments(mediaID string, limit int) ([]models.Comment, error)
	ListReplies(commentID string, limit int) ([]models.Comment, error)

	// Insights
	GetAccountInsights(userID, period string) ([]models.AccountInsight, error)
	GetAudienceDemographics(userID string) ([]models.AudienceDemographic, error)

	// Discovery
	DiscoverUser(userID, username string) (*models.BusinessDiscovery, error)

	// Auth
	ExchangeCodeForToken(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error)
	ExchangeForLongLivedToken(appID, appSecret, shortToken string) (*models.TokenResponse, error)
	RefreshLongLivedToken(token string) (*models.TokenResponse, error)
	GetUserProfile(token string) (*models.UserProfile, error)
}

// GraphClient implements Client using HTTP calls to the Instagram Graph API.
type GraphClient struct {
	httpClient *http.Client
	token      string
	baseURL    string // overridable for testing; defaults to BaseURL
	fbBaseURL  string // overridable for testing; defaults to FacebookBaseURL
}

// NewGraphClient creates a new GraphClient with the given access token.
func NewGraphClient(token string) *GraphClient {
	return &GraphClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		baseURL:    BaseURL,
		fbBaseURL:  FacebookBaseURL,
	}
}

// RateLimitInfo holds rate limit information from response headers.
type RateLimitInfo struct {
	Usage     int
	Limit     int
	Remaining int
}

// doRequest performs an HTTP request and returns the response body.
func (c *GraphClient) doRequest(method, rawURL string) ([]byte, *RateLimitInfo, error) {
	slog.Debug("API request", "method", method, "url", rawURL)

	req, err := http.NewRequest(method, rawURL, http.NoBody)
	if err != nil {
		return nil, nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("reading response: %w", err)
	}

	// Parse rate limit headers
	rateLimit := parseRateLimitHeaders(resp.Header)
	if rateLimit != nil {
		slog.Debug("Rate limit", "usage", rateLimit.Usage, "limit", rateLimit.Limit)
		if rateLimit.Limit > 0 {
			pct := float64(rateLimit.Usage) / float64(rateLimit.Limit) * 100
			if pct >= 80 {
				slog.Warn("Approaching rate limit", "usage_percent", fmt.Sprintf("%.0f%%", pct))
			}
		}
	}

	slog.Debug("API response", "status", resp.StatusCode, "body_length", len(body))

	if resp.StatusCode >= 400 {
		var errResp models.APIErrorResponse
		if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, rateLimit, &APIError{
				StatusCode: resp.StatusCode,
				Message:    errResp.Error.Message,
				Type:       errResp.Error.Type,
				Code:       errResp.Error.Code,
				FBTraceID:  errResp.Error.FBTraceID,
			}
		}
		return nil, rateLimit, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	return body, rateLimit, nil
}

// buildURL constructs a URL with the access token and additional params.
func (c *GraphClient) buildURL(path string, params url.Values) string {
	if params == nil {
		params = url.Values{}
	}
	if c.token != "" {
		params.Set("access_token", c.token)
	}
	return fmt.Sprintf("%s%s?%s", c.baseURL, path, params.Encode())
}

// buildFacebookURL constructs a Facebook Graph API URL with params.
func (c *GraphClient) buildFacebookURL(path string, params url.Values) string {
	if params == nil {
		params = url.Values{}
	}
	base := c.fbBaseURL
	if base == "" {
		base = FacebookBaseURL
	}
	return fmt.Sprintf("%s%s?%s", base, path, params.Encode())
}

// autoPaginate fetches all pages of a paginated response up to the given limit.
func autoPaginate[T any](c *GraphClient, initialURL string, limit int) ([]T, error) {
	var all []T
	nextURL := initialURL

	for nextURL != "" {
		body, _, err := c.doRequest("GET", nextURL)
		if err != nil {
			return nil, err
		}

		var page models.PaginatedResponse[T]
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("parsing paginated response: %w", err)
		}

		all = append(all, page.Data...)

		if limit > 0 && len(all) >= limit {
			all = all[:limit]
			break
		}

		nextURL = page.Paging.Next
	}

	return all, nil
}

func parseRateLimitHeaders(h http.Header) *RateLimitInfo {
	usageStr := h.Get("X-App-Usage")
	if usageStr == "" {
		return nil
	}

	// X-App-Usage is a JSON object like {"call_count":28,"total_cputime":2,"total_time":5}
	var usage struct {
		CallCount int `json:"call_count"`
	}
	if err := json.Unmarshal([]byte(usageStr), &usage); err != nil {
		return nil
	}

	info := &RateLimitInfo{
		Usage: usage.CallCount,
		Limit: 100, // Meta's default limit is 100%
	}

	if remaining := h.Get("X-RateLimit-Remaining"); remaining != "" {
		if v, err := strconv.Atoi(remaining); err == nil {
			info.Remaining = v
		}
	}

	return info
}
