package readiness

import "github.com/prepio/prepio/config"

// CompanyStats holds V1 aggregate answer performance for one company tag.
type CompanyStats struct {
	Company  string
	Answered int
	Correct  int
	ScoreAvg int
}

// CompanyScore is a V1 readiness score for one company.
type CompanyScore struct {
	Company string
	Score   int
}

// ComputeV1CompanyScore derives legacy readiness from tag-based answer history.
func ComputeV1CompanyScore(stats CompanyStats) int {
	if stats.Answered <= 0 {
		return 0
	}
	score := (stats.Correct * 100) / stats.Answered
	if stats.ScoreAvg > 0 {
		score = (score + stats.ScoreAvg) / 2
	}
	if score > config.MaxCompanyReadiness {
		score = config.MaxCompanyReadiness
	}
	return score
}

// ComputeV1Readiness returns V1 company readiness for each target company.
func ComputeV1Readiness(targets []string, byCompany map[string]CompanyStats) []CompanyScore {
	if len(targets) == 0 {
		return []CompanyScore{}
	}

	cards := make([]CompanyScore, 0, len(targets))
	for _, company := range targets {
		row, ok := byCompany[company]
		score := 0
		if ok {
			score = ComputeV1CompanyScore(row)
		}
		cards = append(cards, CompanyScore{Company: company, Score: score})
	}
	return cards
}

// ComputeV1Overall averages V1 company scores for target companies.
func ComputeV1Overall(scores []CompanyScore) int {
	if len(scores) == 0 {
		return 0
	}
	total := 0
	for _, score := range scores {
		total += score.Score
	}
	return total / len(scores)
}

// V1FormulaDescription documents the legacy readiness formula.
const V1FormulaDescription = "(correct_rate * 100 + avg_score) / 2, capped at 95"

// V2FormulaDescription documents the skill-based readiness formula.
const V2FormulaDescription = "sum(skill_mastery * company_weight) / sum(company_weight), capped at 95"
