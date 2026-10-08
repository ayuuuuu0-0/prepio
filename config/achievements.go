package config

// Achievement kinds: each is evaluated from facts the Progress service owns.
const (
	AchievementLessons      = "lessons"       // lessons completed (first completions)
	AchievementStreak       = "streak"        // current streak, in days
	AchievementLevel        = "level"         // current level
	AchievementTopicMastery = "topic_mastery" // best topic readiness, 0-100
	AchievementPromotion    = "promotion"     // promoted in a weekly league at least once
)

// Achievement is one collectible milestone. Order in Achievements is display order.
type Achievement struct {
	Slug        string
	Name        string
	Description string // how to earn it, shown while it is still locked
	Kind        string
	Threshold   int
}

// Achievements is the catalog. Slugs are stable: they are stored per user, so never rename one.
var Achievements = []Achievement{
	{Slug: "first-steps", Name: "First Steps", Description: "Finish your first lesson.", Kind: AchievementLessons, Threshold: 1},
	{Slug: "warming-up", Name: "Warming Up", Description: "Keep a 3-day streak.", Kind: AchievementStreak, Threshold: 3},
	{Slug: "ten-lessons", Name: "Getting Serious", Description: "Finish 10 lessons.", Kind: AchievementLessons, Threshold: 10},
	{Slug: "full-week", Name: "Full Week", Description: "Keep a 7-day streak.", Kind: AchievementStreak, Threshold: 7},
	{Slug: "level-five", Name: "Level Five", Description: "Reach level 5.", Kind: AchievementLevel, Threshold: 5},
	{Slug: "halfway-there", Name: "Halfway There", Description: "Reach 50 readiness in any topic.", Kind: AchievementTopicMastery, Threshold: 50},
	{Slug: "moving-up", Name: "Moving Up", Description: "Get promoted in a weekly league.", Kind: AchievementPromotion, Threshold: 1},
	{Slug: "path-finder", Name: "Path Finder", Description: "Finish 30 lessons.", Kind: AchievementLessons, Threshold: 30},
}

// AchievementFacts are what Progress knows about a learner when it evaluates achievements.
// A fact the caller does not know is left at zero, which never earns anything.
type AchievementFacts struct {
	Lessons          int
	Streak           int
	Level            int
	BestTopicMastery int
	Promoted         bool
}

// EarnedAchievements returns the slugs of every achievement the facts satisfy, in catalog
// order. Callers store them idempotently, so already-earned slugs are harmless.
func EarnedAchievements(f AchievementFacts) []string {
	var earned []string
	for _, a := range Achievements {
		var value int
		switch a.Kind {
		case AchievementLessons:
			value = f.Lessons
		case AchievementStreak:
			value = f.Streak
		case AchievementLevel:
			value = f.Level
		case AchievementTopicMastery:
			value = f.BestTopicMastery
		case AchievementPromotion:
			if f.Promoted {
				value = 1
			}
		}
		if value >= a.Threshold {
			earned = append(earned, a.Slug)
		}
	}
	return earned
}

// AchievementBySlug looks up a catalog entry.
func AchievementBySlug(slug string) (Achievement, bool) {
	for _, a := range Achievements {
		if a.Slug == slug {
			return a, true
		}
	}
	return Achievement{}, false
}
