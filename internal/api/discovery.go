package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pj4533/ig-cli/internal/models"
)

// DiscoverUser looks up a public Business/Creator account via business discovery.
func (c *GraphClient) DiscoverUser(userID, username string) (*models.BusinessDiscovery, error) {
	params := url.Values{}
	fieldStr := fmt.Sprintf("business_discovery.fields(id,username,name,biography,followers_count,media_count,profile_picture_url,website){%s}", username)
	params.Set("fields", fieldStr)

	apiURL := c.buildURL(fmt.Sprintf("/%s", userID), params)

	body, _, err := c.doRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("discovering user %q: %w", username, err)
	}

	var result struct {
		BusinessDiscovery models.BusinessDiscovery `json:"business_discovery"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing discovery response: %w", err)
	}

	return &result.BusinessDiscovery, nil
}
