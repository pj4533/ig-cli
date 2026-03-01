package api

import "github.com/pj4533/ig-cli/internal/models"

// MockClient implements Client for testing.
type MockClient struct {
	ListMediaFn                 func(userID string, limit int) ([]models.Media, error)
	GetMediaInsightsFn          func(mediaID string) ([]models.MediaInsight, error)
	ListCommentsFn              func(mediaID string, limit int) ([]models.Comment, error)
	ListRepliesFn               func(commentID string, limit int) ([]models.Comment, error)
	GetAccountInsightsFn        func(userID string, period string) ([]models.AccountInsight, error)
	GetAudienceDemographicsFn   func(userID string) ([]models.AudienceDemographic, error)
	DiscoverUserFn              func(userID string, username string) (*models.BusinessDiscovery, error)
	ExchangeCodeForTokenFn      func(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error)
	ExchangeForLongLivedTokenFn func(appID, appSecret, shortToken string) (*models.TokenResponse, error)
	RefreshLongLivedTokenFn     func(token string) (*models.TokenResponse, error)
	GetUserProfileFn            func(token string) (*models.UserProfile, error)
}

func (m *MockClient) ListMedia(userID string, limit int) ([]models.Media, error) {
	if m.ListMediaFn != nil {
		return m.ListMediaFn(userID, limit)
	}
	return nil, nil
}

func (m *MockClient) GetMediaInsights(mediaID string) ([]models.MediaInsight, error) {
	if m.GetMediaInsightsFn != nil {
		return m.GetMediaInsightsFn(mediaID)
	}
	return nil, nil
}

func (m *MockClient) ListComments(mediaID string, limit int) ([]models.Comment, error) {
	if m.ListCommentsFn != nil {
		return m.ListCommentsFn(mediaID, limit)
	}
	return nil, nil
}

func (m *MockClient) ListReplies(commentID string, limit int) ([]models.Comment, error) {
	if m.ListRepliesFn != nil {
		return m.ListRepliesFn(commentID, limit)
	}
	return nil, nil
}

func (m *MockClient) GetAccountInsights(userID, period string) ([]models.AccountInsight, error) {
	if m.GetAccountInsightsFn != nil {
		return m.GetAccountInsightsFn(userID, period)
	}
	return nil, nil
}

func (m *MockClient) GetAudienceDemographics(userID string) ([]models.AudienceDemographic, error) {
	if m.GetAudienceDemographicsFn != nil {
		return m.GetAudienceDemographicsFn(userID)
	}
	return nil, nil
}

func (m *MockClient) DiscoverUser(userID, username string) (*models.BusinessDiscovery, error) {
	if m.DiscoverUserFn != nil {
		return m.DiscoverUserFn(userID, username)
	}
	return nil, nil
}

func (m *MockClient) ExchangeCodeForToken(appID, appSecret, redirectURI, code string) (*models.TokenResponse, error) {
	if m.ExchangeCodeForTokenFn != nil {
		return m.ExchangeCodeForTokenFn(appID, appSecret, redirectURI, code)
	}
	return nil, nil
}

func (m *MockClient) ExchangeForLongLivedToken(appID, appSecret, shortToken string) (*models.TokenResponse, error) {
	if m.ExchangeForLongLivedTokenFn != nil {
		return m.ExchangeForLongLivedTokenFn(appID, appSecret, shortToken)
	}
	return nil, nil
}

func (m *MockClient) RefreshLongLivedToken(token string) (*models.TokenResponse, error) {
	if m.RefreshLongLivedTokenFn != nil {
		return m.RefreshLongLivedTokenFn(token)
	}
	return nil, nil
}

func (m *MockClient) GetUserProfile(token string) (*models.UserProfile, error) {
	if m.GetUserProfileFn != nil {
		return m.GetUserProfileFn(token)
	}
	return nil, nil
}
