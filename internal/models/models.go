package models

import "time"

// TokenResponse represents the OAuth token exchange response.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// UserProfile represents an Instagram user profile.
type UserProfile struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name,omitempty"`
}

// Media represents an Instagram media object.
type Media struct {
	ID            string    `json:"id"`
	Caption       string    `json:"caption,omitempty"`
	MediaType     string    `json:"media_type"`
	MediaURL      string    `json:"media_url,omitempty"`
	Permalink     string    `json:"permalink,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
	LikeCount     int       `json:"like_count"`
	CommentsCount int       `json:"comments_count"`
}

// Comment represents a comment on a media object.
type Comment struct {
	ID        string                      `json:"id"`
	Text      string                      `json:"text"`
	Username  string                      `json:"username"`
	Timestamp time.Time                   `json:"timestamp"`
	LikeCount int                         `json:"like_count"`
	Replies   *PaginatedResponse[Comment] `json:"replies,omitempty"`
}

// MetricValue represents a single metric data point.
type MetricValue struct {
	Value   interface{} `json:"value"`
	EndTime string      `json:"end_time,omitempty"`
}

// MediaInsight represents insight metrics for a media object.
type MediaInsight struct {
	Name        string        `json:"name"`
	Period      string        `json:"period"`
	Values      []MetricValue `json:"values"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ID          string        `json:"id"`
}

// AccountInsight represents account-level insight metrics.
type AccountInsight struct {
	Name        string        `json:"name"`
	Period      string        `json:"period"`
	Values      []MetricValue `json:"values"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	ID          string        `json:"id"`
}

// AudienceDemographic represents audience demographic data.
type AudienceDemographic struct {
	Name        string        `json:"name"`
	Period      string        `json:"period"`
	Values      []MetricValue `json:"values"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
}

// BusinessDiscovery represents discovered business account data.
type BusinessDiscovery struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name,omitempty"`
	Biography      string `json:"biography,omitempty"`
	FollowersCount int    `json:"followers_count"`
	MediaCount     int    `json:"media_count"`
	ProfilePicURL  string `json:"profile_picture_url,omitempty"`
	Website        string `json:"website,omitempty"`
}

// Paging holds pagination cursors.
type Paging struct {
	Cursors struct {
		Before string `json:"before"`
		After  string `json:"after"`
	} `json:"cursors"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
}

// PaginatedResponse is a generic paginated API response.
type PaginatedResponse[T any] struct {
	Data   []T    `json:"data"`
	Paging Paging `json:"paging"`
}

// APIErrorResponse represents an error from the Instagram Graph API.
type APIErrorResponse struct {
	Error struct {
		Message   string `json:"message"`
		Type      string `json:"type"`
		Code      int    `json:"code"`
		FBTraceID string `json:"fbtrace_id"`
	} `json:"error"`
}
