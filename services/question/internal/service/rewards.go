package service

import (
	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/store"
)

// computeRewards estimates XP and gems from difficulty, score, and company tags.
func computeRewards(question *store.Question, eval EvaluationResult) (xp, gems int) {
	if !eval.Correct || eval.Score < constants.MinEvaluationScore {
		return 0, 0
	}

	xpBase, ok := config.XPByDifficulty[question.Difficulty]
	if !ok {
		xpBase = config.XPByDifficulty["medium"]
	}
	gemsBase, ok := config.GemsByDifficulty[question.Difficulty]
	if !ok {
		gemsBase = config.GemsByDifficulty["medium"]
	}

	for _, tag := range question.CompanyTags {
		if config.TopTierCompanies[tag] {
			xpBase = int(float64(xpBase) * config.TopTierCompanyXPMultiplier)
			break
		}
	}

	xp = xpBase * eval.Score / 100
	if xp < 1 && eval.Correct {
		xp = 1
	}

	gems = gemsBase
	if eval.Score < 80 {
		gems = gemsBase / 2
	}
	if gems < 1 && eval.Correct {
		gems = 1
	}

	return xp, gems
}
