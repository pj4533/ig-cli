package api

import (
	"testing"

	"github.com/pj4533/ig-cli/internal/models"
)

func TestGetAccountInsights(t *testing.T) {
	mock := &MockClient{
		GetAccountInsightsFn: func(userID string, period string) ([]models.AccountInsight, error) {
			if userID != "12345" {
				t.Errorf("userID = %q, want %q", userID, "12345")
			}
			if period != "day" {
				t.Errorf("period = %q, want %q", period, "day")
			}
			return []models.AccountInsight{
				{Name: "impressions", Period: "day", Values: []models.MetricValue{{Value: 5000}}},
				{Name: "reach", Period: "day", Values: []models.MetricValue{{Value: 3000}}},
			}, nil
		},
	}

	insights, err := mock.GetAccountInsights("12345", "day")
	if err != nil {
		t.Fatalf("GetAccountInsights error: %v", err)
	}
	if len(insights) != 2 {
		t.Fatalf("insights length = %d, want 2", len(insights))
	}
	if insights[0].Name != "impressions" {
		t.Errorf("Name = %q, want %q", insights[0].Name, "impressions")
	}
}

func TestGetAudienceDemographics(t *testing.T) {
	mock := &MockClient{
		GetAudienceDemographicsFn: func(userID string) ([]models.AudienceDemographic, error) {
			if userID != "12345" {
				t.Errorf("userID = %q, want %q", userID, "12345")
			}
			return []models.AudienceDemographic{
				{Name: "audience_city", Period: "lifetime"},
				{Name: "audience_country", Period: "lifetime"},
			}, nil
		},
	}

	demographics, err := mock.GetAudienceDemographics("12345")
	if err != nil {
		t.Fatalf("GetAudienceDemographics error: %v", err)
	}
	if len(demographics) != 2 {
		t.Fatalf("demographics length = %d, want 2", len(demographics))
	}
}
