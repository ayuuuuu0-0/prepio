package config

import (
	"reflect"
	"testing"
)

func TestEarnedAchievements(t *testing.T) {
	cases := []struct {
		name  string
		facts AchievementFacts
		want  []string
	}{
		{"nothing yet", AchievementFacts{}, nil},
		{"first lesson", AchievementFacts{Lessons: 1, Level: 1}, []string{"first-steps"}},
		{"streak alone (lesson facts unknown)", AchievementFacts{Streak: 7}, []string{"warming-up", "full-week"}},
		{"thresholds are inclusive", AchievementFacts{Lessons: 10, Level: 5, BestTopicMastery: 50}, []string{"first-steps", "ten-lessons", "level-five", "halfway-there"}},
		{"just below every threshold", AchievementFacts{Lessons: 0, Streak: 2, Level: 4, BestTopicMastery: 49}, nil},
		{"promotion", AchievementFacts{Promoted: true}, []string{"moving-up"}},
		{"everything", AchievementFacts{Lessons: 40, Streak: 9, Level: 8, BestTopicMastery: 80, Promoted: true},
			[]string{"first-steps", "warming-up", "ten-lessons", "full-week", "level-five", "halfway-there", "moving-up", "path-finder"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := EarnedAchievements(c.facts); !reflect.DeepEqual(got, c.want) {
				t.Fatalf("EarnedAchievements(%+v) = %v, want %v", c.facts, got, c.want)
			}
		})
	}
}

func TestAchievementCatalogIsWellFormed(t *testing.T) {
	seen := map[string]bool{}
	kinds := map[string]bool{AchievementLessons: true, AchievementStreak: true, AchievementLevel: true, AchievementTopicMastery: true, AchievementPromotion: true}
	for _, a := range Achievements {
		if a.Slug == "" || a.Name == "" || a.Description == "" || a.Threshold < 1 || !kinds[a.Kind] {
			t.Fatalf("malformed achievement %+v", a)
		}
		if seen[a.Slug] {
			t.Fatalf("duplicate slug %s", a.Slug)
		}
		seen[a.Slug] = true
		if got, ok := AchievementBySlug(a.Slug); !ok || got != a {
			t.Fatalf("AchievementBySlug(%s) failed", a.Slug)
		}
	}
	if _, ok := AchievementBySlug("nope"); ok {
		t.Fatal("unknown slug found")
	}
}
