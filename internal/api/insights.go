package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pj4533/ig-cli/internal/models"
)

// GetAccountInsights fetches account-level insights.
func (c *GraphClient) GetAccountInsights(userID, period string) ([]models.AccountInsight, error) {
	if period == "" {
		period = "day"
	}

	params := url.Values{
		"metric": {"impressions,reach,profile_views,website_clicks,follower_count,email_contacts,phone_call_clicks,text_message_clicks,get_directions_clicks"},
		"period": {period},
	}

	apiURL := c.buildURL(fmt.Sprintf("/%s/insights", userID), params)

	body, _, err := c.doRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("getting account insights: %w", err)
	}

	var resp models.PaginatedResponse[models.AccountInsight]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing account insights: %w", err)
	}

	return resp.Data, nil
}

// GetAudienceDemographics fetches audience demographic data.
func (c *GraphClient) GetAudienceDemographics(userID string) ([]models.AudienceDemographic, error) {
	params := url.Values{
		"metric": {"audience_city,audience_country,audience_gender_age,audience_locale"},
		"period": {"lifetime"},
	}

	apiURL := c.buildURL(fmt.Sprintf("/%s/insights", userID), params)

	body, _, err := c.doRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("getting audience demographics: %w", err)
	}

	var resp models.PaginatedResponse[models.AudienceDemographic]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing audience demographics: %w", err)
	}

	return resp.Data, nil
}
