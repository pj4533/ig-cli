package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestMediaJSON(t *testing.T) {
	jsonStr := `{
		"id": "123",
		"caption": "Test post",
		"media_type": "IMAGE",
		"media_url": "https://example.com/img.jpg",
		"permalink": "https://instagram.com/p/abc",
		"timestamp": "2025-01-15T10:30:00+00:00",
		"like_count": 42,
		"comments_count": 5
	}`

	var m Media
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if m.ID != "123" {
		t.Errorf("ID = %q, want %q", m.ID, "123")
	}
	if m.Caption != "Test post" {
		t.Errorf("Caption = %q, want %q", m.Caption, "Test post")
	}
	if m.MediaType != "IMAGE" {
		t.Errorf("MediaType = %q, want %q", m.MediaType, "IMAGE")
	}
	if m.LikeCount != 42 {
		t.Errorf("LikeCount = %d, want %d", m.LikeCount, 42)
	}
	if m.CommentsCount != 5 {
		t.Errorf("CommentsCount = %d, want %d", m.CommentsCount, 5)
	}
}

func TestMediaMarshal(t *testing.T) {
	m := Media{
		ID:            "456",
		Caption:       "Hello world",
		MediaType:     "VIDEO",
		Timestamp:     time.Date(2025, 6, 15, 12, 0, 0, 0, time.UTC),
		LikeCount:     100,
		CommentsCount: 10,
	}

	data, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded Media
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.ID != m.ID {
		t.Errorf("ID = %q, want %q", decoded.ID, m.ID)
	}
	if decoded.LikeCount != m.LikeCount {
		t.Errorf("LikeCount = %d, want %d", decoded.LikeCount, m.LikeCount)
	}
}

func TestCommentJSON(t *testing.T) {
	jsonStr := `{
		"id": "c1",
		"text": "Great post!",
		"username": "user1",
		"timestamp": "2025-01-15T12:00:00+00:00",
		"like_count": 3
	}`

	var c Comment
	if err := json.Unmarshal([]byte(jsonStr), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if c.ID != "c1" {
		t.Errorf("ID = %q, want %q", c.ID, "c1")
	}
	if c.Text != "Great post!" {
		t.Errorf("Text = %q, want %q", c.Text, "Great post!")
	}
	if c.Username != "user1" {
		t.Errorf("Username = %q, want %q", c.Username, "user1")
	}
}

func TestTokenResponseJSON(t *testing.T) {
	jsonStr := `{
		"access_token": "abc123",
		"token_type": "bearer",
		"expires_in": 5184000
	}`

	var tr TokenResponse
	if err := json.Unmarshal([]byte(jsonStr), &tr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if tr.AccessToken != "abc123" {
		t.Errorf("AccessToken = %q, want %q", tr.AccessToken, "abc123")
	}
	if tr.ExpiresIn != 5184000 {
		t.Errorf("ExpiresIn = %d, want %d", tr.ExpiresIn, 5184000)
	}
}

func TestUserProfileJSON(t *testing.T) {
	jsonStr := `{"id": "12345", "username": "testuser", "name": "Test User"}`

	var p UserProfile
	if err := json.Unmarshal([]byte(jsonStr), &p); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if p.ID != "12345" {
		t.Errorf("ID = %q, want %q", p.ID, "12345")
	}
	if p.Username != "testuser" {
		t.Errorf("Username = %q, want %q", p.Username, "testuser")
	}
}

func TestBusinessDiscoveryJSON(t *testing.T) {
	jsonStr := `{
		"id": "bd1",
		"username": "business",
		"name": "Business Account",
		"biography": "A business",
		"followers_count": 10000,
		"media_count": 500
	}`

	var bd BusinessDiscovery
	if err := json.Unmarshal([]byte(jsonStr), &bd); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if bd.FollowersCount != 10000 {
		t.Errorf("FollowersCount = %d, want %d", bd.FollowersCount, 10000)
	}
	if bd.MediaCount != 500 {
		t.Errorf("MediaCount = %d, want %d", bd.MediaCount, 500)
	}
}

func TestPaginatedResponseJSON(t *testing.T) {
	jsonStr := `{
		"data": [{"id": "1"}, {"id": "2"}],
		"paging": {
			"cursors": {"before": "abc", "after": "def"},
			"next": "https://example.com/next"
		}
	}`

	var pr PaginatedResponse[Media]
	if err := json.Unmarshal([]byte(jsonStr), &pr); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if len(pr.Data) != 2 {
		t.Errorf("Data length = %d, want %d", len(pr.Data), 2)
	}
	if pr.Paging.Next != "https://example.com/next" {
		t.Errorf("Next = %q, want %q", pr.Paging.Next, "https://example.com/next")
	}
}

func TestMediaInsightJSON(t *testing.T) {
	jsonStr := `{
		"name": "impressions",
		"period": "lifetime",
		"values": [{"value": 1234}],
		"title": "Impressions",
		"description": "Total impressions",
		"id": "123/insights/impressions/lifetime"
	}`

	var mi MediaInsight
	if err := json.Unmarshal([]byte(jsonStr), &mi); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if mi.Name != "impressions" {
		t.Errorf("Name = %q, want %q", mi.Name, "impressions")
	}
	if len(mi.Values) != 1 {
		t.Errorf("Values length = %d, want %d", len(mi.Values), 1)
	}
}

func TestAudienceDemographicJSON(t *testing.T) {
	jsonStr := `{
		"name": "audience_city",
		"period": "lifetime",
		"values": [{"value": {"New York": 500, "Los Angeles": 300}}],
		"title": "Audience City",
		"description": "Cities"
	}`

	var ad AudienceDemographic
	if err := json.Unmarshal([]byte(jsonStr), &ad); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if ad.Name != "audience_city" {
		t.Errorf("Name = %q, want %q", ad.Name, "audience_city")
	}
}

func TestAPIErrorResponseJSON(t *testing.T) {
	jsonStr := `{
		"error": {
			"message": "Invalid token",
			"type": "OAuthException",
			"code": 190,
			"fbtrace_id": "abc123"
		}
	}`

	var aer APIErrorResponse
	if err := json.Unmarshal([]byte(jsonStr), &aer); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if aer.Error.Code != 190 {
		t.Errorf("Code = %d, want %d", aer.Error.Code, 190)
	}
	if aer.Error.Message != "Invalid token" {
		t.Errorf("Message = %q, want %q", aer.Error.Message, "Invalid token")
	}
}
