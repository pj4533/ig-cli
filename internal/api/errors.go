package api

import (
	"fmt"
)

// APIError represents a structured error from the Instagram Graph API.
type APIError struct {
	StatusCode int
	Message    string
	Type       string
	Code       int
	FBTraceID  string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Instagram API error %d (code %d): %s", e.StatusCode, e.Code, e.Message)
}

// IsRateLimited returns true if the error is a rate limit error.
func (e *APIError) IsRateLimited() bool {
	return e.Code == 4 || e.Code == 32 || e.Code == 613
}

// IsAuthExpired returns true if the error indicates an expired or invalid token.
func (e *APIError) IsAuthExpired() bool {
	return e.Code == 190 || e.StatusCode == 401
}
