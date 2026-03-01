package api

import (
	"fmt"
	"net/url"

	"github.com/pj4533/ig-cli/internal/models"
)

// ListComments fetches comments on a media object.
func (c *GraphClient) ListComments(mediaID string, limit int) ([]models.Comment, error) {
	params := url.Values{
		"fields": {"id,text,username,timestamp,like_count"},
	}
	if limit > 0 && limit <= 100 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	apiURL := c.buildURL(fmt.Sprintf("/%s/comments", mediaID), params)
	return autoPaginate[models.Comment](c, apiURL, limit)
}

// ListReplies fetches replies to a comment.
func (c *GraphClient) ListReplies(commentID string, limit int) ([]models.Comment, error) {
	params := url.Values{
		"fields": {"id,text,username,timestamp,like_count"},
	}
	if limit > 0 && limit <= 100 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	apiURL := c.buildURL(fmt.Sprintf("/%s/replies", commentID), params)
	return autoPaginate[models.Comment](c, apiURL, limit)
}
