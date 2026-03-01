package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pj4533/ig-cli/internal/models"
)

// ListMedia fetches media objects for a user.
func (c *GraphClient) ListMedia(userID string, limit int) ([]models.Media, error) {
	params := url.Values{
		"fields": {"id,caption,media_type,media_url,permalink,timestamp,like_count,comments_count"},
	}
	if limit > 0 && limit <= 100 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	apiURL := c.buildURL(fmt.Sprintf("/%s/media", userID), params)
	return autoPaginate[models.Media](c, apiURL, limit)
}

// GetMediaInsights fetches insight metrics for a media object.
func (c *GraphClient) GetMediaInsights(mediaID string) ([]models.MediaInsight, error) {
	params := url.Values{
		"metric": {"impressions,reach,engagement,saved,video_views,likes,comments,shares"},
	}

	apiURL := c.buildURL(fmt.Sprintf("/%s/insights", mediaID), params)

	body, _, err := c.doRequest("GET", apiURL)
	if err != nil {
		return nil, fmt.Errorf("getting media insights: %w", err)
	}

	var resp models.PaginatedResponse[models.MediaInsight]
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parsing media insights: %w", err)
	}

	return resp.Data, nil
}
