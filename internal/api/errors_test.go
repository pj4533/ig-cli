package api

import (
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	e := &APIError{
		StatusCode: 400,
		Code:       190,
		Message:    "Invalid token",
	}

	got := e.Error()
	want := "Instagram API error 400 (code 190): Invalid token"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAPIError_IsRateLimited(t *testing.T) {
	tests := []struct {
		code int
		want bool
	}{
		{4, true},
		{32, true},
		{613, true},
		{190, false},
		{0, false},
	}

	for _, tt := range tests {
		e := &APIError{Code: tt.code}
		if got := e.IsRateLimited(); got != tt.want {
			t.Errorf("IsRateLimited() for code %d = %v, want %v", tt.code, got, tt.want)
		}
	}
}

func TestAPIError_IsAuthExpired(t *testing.T) {
	tests := []struct {
		code       int
		statusCode int
		want       bool
	}{
		{190, 400, true},
		{0, 401, true},
		{4, 429, false},
		{0, 200, false},
	}

	for _, tt := range tests {
		e := &APIError{Code: tt.code, StatusCode: tt.statusCode}
		if got := e.IsAuthExpired(); got != tt.want {
			t.Errorf("IsAuthExpired() for code=%d status=%d = %v, want %v", tt.code, tt.statusCode, got, tt.want)
		}
	}
}
