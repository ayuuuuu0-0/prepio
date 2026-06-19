package constants

// Readiness score source values stored on user_skill_scores.
const (
	ReadinessSourceLive     = "live"
	ReadinessSourceBackfill = "backfill"
)

// Supported company slugs with seeded skill weight profiles.
var ReadinessCompanies = []string{
	"google",
	"amazon",
	"meta",
	"uber",
}
