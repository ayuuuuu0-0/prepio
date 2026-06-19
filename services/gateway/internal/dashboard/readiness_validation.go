package dashboard

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/shared/readiness"
)

// ReadinessValidationResponse exposes V1 and V2 readiness side by side for validation.
type ReadinessValidationResponse = readiness.ValidationResponse

// GetReadinessValidation returns parallel V1/V2 readiness with explanations and gaps.
func (s *Service) GetReadinessValidation(ctx context.Context, token string) (*ReadinessValidationResponse, error) {
	profile, err := s.fetchProfile(ctx, token)
	if err != nil {
		return nil, err
	}

	v1Stats, err := s.fetchReadinessStats(ctx, token)
	if err != nil {
		return nil, err
	}

	v1Snapshot := buildV1Snapshot(profile.TargetCompanies, v1Stats)

	v2Dashboard, err := s.fetchReadinessDashboard(ctx, token)
	if err != nil {
		return nil, err
	}

	comparison := buildComparison(profile.TargetCompanies, v1Snapshot, v2Dashboard.CompanyReadiness)

	return &ReadinessValidationResponse{
		ReadinessV2Enabled: config.ReadinessV2Enabled(),
		V1:                 v1Snapshot,
		V2:                 *v2Dashboard,
		Comparison:         comparison,
	}, nil
}

func buildV1Snapshot(targets []string, stats *readinessStatsPayload) readiness.V1Snapshot {
	byCompany := map[string]readiness.CompanyStats{}
	if stats != nil {
		for _, row := range stats.ByCompany {
			byCompany[row.Company] = readiness.CompanyStats{
				Company:  row.Company,
				Answered: row.Answered,
				Correct:  row.Correct,
				ScoreAvg: row.ScoreAvg,
			}
		}
	}

	v1Scores := readiness.ComputeV1Readiness(targets, byCompany)
	cards := make([]readiness.ScoreCard, 0, len(v1Scores))
	for _, score := range v1Scores {
		cards = append(cards, readiness.ScoreCard{
			Company: score.Company,
			Score:   score.Score,
		})
	}

	return readiness.V1Snapshot{
		Companies: cards,
		Overall:   readiness.ComputeV1Overall(v1Scores),
		Version:   "v1",
		Formula:   readiness.V1FormulaDescription,
	}
}

func buildComparison(
	targets []string,
	v1 readiness.V1Snapshot,
	v2 readiness.CompanyResponse,
) readiness.ComparisonSnapshot {
	v2ByCompany := map[string]int{}
	for _, entry := range v2.Companies {
		v2ByCompany[entry.Company] = entry.Readiness
	}

	v1ByCompany := map[string]int{}
	for _, entry := range v1.Companies {
		v1ByCompany[entry.Company] = entry.Score
	}

	comparisons := make([]readiness.CompanyComparison, 0, len(targets))
	for _, company := range targets {
		v1Score := v1ByCompany[company]
		v2Score := v2ByCompany[company]
		comparisons = append(comparisons, readiness.CompanyComparison{
			Company: company,
			V1Score: v1Score,
			V2Score: v2Score,
			Delta:   v2Score - v1Score,
		})
	}

	return readiness.ComparisonSnapshot{
		OverallV1:    v1.Overall,
		OverallV2:    v2.Overall,
		OverallDelta: v2.Overall - v1.Overall,
		ByCompany:    comparisons,
		V1Formula:    readiness.V1FormulaDescription,
		V2Formula:    readiness.V2FormulaDescription,
	}
}

func (s *Service) fetchReadinessDashboard(ctx context.Context, token string) (*readiness.DashboardResponse, error) {
	if len(s.progressURL) == 0 {
		return &readiness.DashboardResponse{}, nil
	}
	body, err := s.get(ctx, s.progressURL+"/api/v1/readiness/dashboard", token)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data readiness.DashboardResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("decode readiness dashboard: %w", err)
	}
	return &envelope.Data, nil
}

func (s *Service) fetchCompanyReadinessV2(ctx context.Context, token string) ([]ReadinessCard, error) {
	if len(s.progressURL) == 0 {
		return []ReadinessCard{}, nil
	}
	body, err := s.get(ctx, s.progressURL+"/api/v1/companies/readiness", token)
	if err != nil {
		return nil, err
	}
	var envelope struct {
		Data readiness.CompanyResponse `json:"data"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("decode company readiness v2: %w", err)
	}

	cards := make([]ReadinessCard, 0, len(envelope.Data.Companies))
	for _, company := range envelope.Data.Companies {
		cards = append(cards, ReadinessCard{
			Company: company.Company,
			Score:   company.Readiness,
		})
	}
	return cards, nil
}
