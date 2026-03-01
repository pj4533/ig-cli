package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pj4533/ig-cli/internal/models"
)

// ExchangeCodeForToken exchanges an authorization code for a short-lived access token.
func (c *GraphClient) ExchangeCodeForToken(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error) {
	params := url.Values{
		"client_id":     {appID},
		"client_secret": {appSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {redirectURI},
		"code":          {code},
	}

	base := "https://api.instagram.com"
	if c.baseURL != "" && c.baseURL != BaseURL {
		base = c.baseURL
	}
	apiURL := fmt.Sprintf("%s/oauth/access_token?%s", base, params.Encode())

	body, _, err := c.doRequest("POST", apiURL)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %w", err)
	}

	var resp models.TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}

	return &resp, nil
}

// ExchangeForLongLivedToken exchanges a short-lived token for a long-lived one (60 days).
func (c *GraphClient) ExchangeForLongLivedToken(appID, appSecret, shortToken string) (*models.TokenResponse, error) {
	params := url.Values{
		"grant_type":    {"ig_exchange_token"},
		"client_secret": {appSecret},
		"access_token":  {shortToken},
	}

	apiURL := c.buildFacebookURL("/oauth/access_token", params)

	body, _, err := c.doRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("exchanging for long-lived token: %w", err)
	}

	var resp models.TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing long-lived token response: %w", err)
	}

	return &resp, nil
}

// RefreshLongLivedToken refreshes a long-lived token.
func (c *GraphClient) RefreshLongLivedToken(token string) (*models.TokenResponse, error) {
	params := url.Values{
		"grant_type":   {"ig_refresh_token"},
		"access_token": {token},
	}

	apiURL := c.buildFacebookURL("/oauth/access_token", params)

	body, _, err := c.doRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("refreshing token: %w", err)
	}

	var resp models.TokenResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing refresh token response: %w", err)
	}

	return &resp, nil
}

// GetUserProfile fetches the user profile using the given token.
func (c *GraphClient) GetUserProfile(token string) (*models.UserProfile, error) {
	params := url.Values{
		"fields":       {"id,username,name"},
		"access_token": {token},
	}

	apiURL := fmt.Sprintf("%s/me?%s", c.baseURL, params.Encode())

	body, _, err := c.doRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("getting user profile: %w", err)
	}

	var profile models.UserProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fmt.Errorf("parsing user profile: %w", err)
	}

	return &profile, nil
}
