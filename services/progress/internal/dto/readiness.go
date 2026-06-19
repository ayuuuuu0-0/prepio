package dto

import sharedreadiness "github.com/prepio/prepio/shared/readiness"

// ReadinessScoreCard is a company readiness score used in V1 and comparison views.
type ReadinessScoreCard = sharedreadiness.ScoreCard

// ReadinessV1Snapshot captures legacy readiness output.
type ReadinessV1Snapshot = sharedreadiness.V1Snapshot

// SkillReadinessEntry is a user's mastery score for one skill.
type SkillReadinessEntry = sharedreadiness.SkillEntry

// SkillSummary highlights a skill for top/weakest lists.
type SkillSummary = sharedreadiness.SkillSummary

// SkillGap identifies a weighted skill holding back company readiness.
type SkillGap = sharedreadiness.SkillGap

// ReadinessExplanation describes why a readiness score looks the way it does.
type ReadinessExplanation = sharedreadiness.Explanation

// SkillReadinessResponse is returned by GET /api/v1/skills/readiness.
type SkillReadinessResponse = sharedreadiness.SkillMasteryResponse

// SkillContribution breaks down one skill's contribution to company readiness.
type SkillContribution = sharedreadiness.SkillContribution

// CompanyReadinessEntry is readiness for one target company.
type CompanyReadinessEntry = sharedreadiness.CompanyEntry

// CompanyReadinessResponse is returned by GET /api/v1/companies/readiness.
type CompanyReadinessResponse = sharedreadiness.CompanyResponse

// ReadinessDashboardResponse aggregates skill and company readiness for validation.
type ReadinessDashboardResponse = sharedreadiness.DashboardResponse

// CompanyComparison compares V1 and V2 readiness for one company.
type CompanyComparison = sharedreadiness.CompanyComparison

// ReadinessComparisonSnapshot holds side-by-side V1 and V2 scores.
type ReadinessComparisonSnapshot = sharedreadiness.ComparisonSnapshot

// ReadinessCompareResponse is returned by internal/compare and dashboard validation APIs.
type ReadinessCompareResponse = sharedreadiness.ValidationResponse
