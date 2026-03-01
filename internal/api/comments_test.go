package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pj4533/ig-cli/internal/models"
)

func TestListComments(t *testing.T) {
	mock := &MockClient{
		ListCommentsFn: func(mediaID string, limit int) ([]models.Comment, error) {
			if mediaID != "m1" {
				t.Errorf("mediaID = %q, want %q", mediaID, "m1")
			}
			return []models.Comment{
				{ID: "c1", Text: "Great!", Username: "user1", LikeCount: 2},
				{ID: "c2", Text: "Nice!", Username: "user2", LikeCount: 0},
			}, nil
		},
	}

	comments, err := mock.ListComments("m1", 0)
	if err != nil {
		t.Fatalf("ListComments error: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("comments length = %d, want 2", len(comments))
	}
	if comments[0].Username != "user1" {
		t.Errorf("Username = %q, want %q", comments[0].Username, "user1")
	}
}

func TestListReplies(t *testing.T) {
	mock := &MockClient{
		ListRepliesFn: func(commentID string, limit int) ([]models.Comment, error) {
			if commentID != "c1" {
				t.Errorf("commentID = %q, want %q", commentID, "c1")
			}
			return []models.Comment{
				{ID: "r1", Text: "Thanks!", Username: "author"},
			}, nil
		},
	}

	replies, err := mock.ListReplies("c1", 0)
	if err != nil {
		t.Fatalf("ListReplies error: %v", err)
	}
	if len(replies) != 1 {
		t.Fatalf("replies length = %d, want 1", len(replies))
	}
}

func TestListComments_HTTPServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := models.PaginatedResponse[models.Comment]{
			Data: []models.Comment{
				{ID: "c1", Text: "Hello", Username: "user1"},
			},
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &GraphClient{
		httpClient: server.Client(),
		token:      "test-token",
	}

	results, err := autoPaginate[models.Comment](client, server.URL+"?access_token=test-token", 0)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("results length = %d, want 1", len(results))
	}
	if results[0].Text != "Hello" {
		t.Errorf("Text = %q, want %q", results[0].Text, "Hello")
	}
}
