package dto

// SkillReadinessEntry is a user's mastery score for one skill.
type SkillReadinessEntry struct {
	SkillSlug       string `json:"skill_slug"`
	SkillName       string `json:"skill_name"`
	Mastery         int    `json:"mastery"`
	Attempts        int    `json:"attempts"`
	LastPracticedAt string `json:"last_practiced_at,omitempty"`
}

// SkillSummary highlights a skill for top/weakest lists.
type SkillSummary struct {
	SkillSlug string `json:"skill_slug"`
	SkillName string `json:"skill_name"`
	Mastery   int    `json:"mastery"`
	Attempts  int    `json:"attempts"`
}

// SkillGap identifies a skill holding back mastery.
type SkillGap struct {
	SkillSlug   string `json:"skill_slug"`
	SkillName   string `json:"skill_name"`
	Mastery     int    `json:"mastery"`
	GapScore    int    `json:"gap_score"`
	Explanation string `json:"explanation"`
}

// ReadinessExplanation describes why a readiness score looks the way it does.
type ReadinessExplanation struct {
	Scope   string   `json:"scope"`
	Summary string   `json:"summary"`
	Details []string `json:"details"`
}

// SkillReadinessResponse is returned by GET /api/v1/skills/readiness.
type SkillReadinessResponse struct {
	Skills        []SkillReadinessEntry  `json:"skills"`
	Overall       int                    `json:"overall"`
	TopSkills     []SkillSummary         `json:"top_skills"`
	WeakestSkills []SkillSummary         `json:"weakest_skills"`
	SkillGaps     []SkillGap             `json:"skill_gaps"`
	Explanations  []ReadinessExplanation `json:"explanations"`
	Version       string                 `json:"version"`
}
