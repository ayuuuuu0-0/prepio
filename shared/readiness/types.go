package readiness

// ScoreCard is a company readiness score.
type ScoreCard struct {
	Company string `json:"company"`
	Score   int    `json:"score"`
}

// V1Snapshot captures legacy readiness output.
type V1Snapshot struct {
	Companies []ScoreCard `json:"companies"`
	Overall   int         `json:"overall"`
	Version   string      `json:"version"`
	Formula   string      `json:"formula"`
}

// SkillEntry is a user's mastery score for one skill.
type SkillEntry struct {
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

// SkillGap identifies a weighted skill holding back company readiness.
type SkillGap struct {
	Company     string `json:"company"`
	SkillSlug   string `json:"skill_slug"`
	SkillName   string `json:"skill_name"`
	Mastery     int    `json:"mastery"`
	Weight      int    `json:"weight"`
	GapScore    int    `json:"gap_score"`
	Explanation string `json:"explanation"`
}

// Explanation describes why a readiness score looks the way it does.
type Explanation struct {
	Scope   string   `json:"scope"`
	Summary string   `json:"summary"`
	Details []string `json:"details"`
}

// SkillMasteryResponse is returned by GET /api/v1/skills/readiness.
type SkillMasteryResponse struct {
	Skills        []SkillEntry  `json:"skills"`
	Overall       int           `json:"overall"`
	TopSkills     []SkillSummary `json:"top_skills"`
	WeakestSkills []SkillSummary `json:"weakest_skills"`
	SkillGaps     []SkillGap    `json:"skill_gaps"`
	Explanations  []Explanation `json:"explanations"`
	Version       string        `json:"version"`
}

// SkillContribution breaks down one skill's contribution to company readiness.
type SkillContribution struct {
	SkillSlug string `json:"skill_slug"`
	SkillName string `json:"skill_name"`
	Mastery   int    `json:"mastery"`
	Weight    int    `json:"weight"`
}

// CompanyEntry is readiness for one target company.
type CompanyEntry struct {
	Company            string              `json:"company"`
	Readiness          int                 `json:"readiness"`
	SkillContributions []SkillContribution `json:"skill_contributions"`
	TopSkills          []SkillSummary      `json:"top_skills"`
	WeakestSkills      []SkillSummary      `json:"weakest_skills"`
	SkillGaps          []SkillGap          `json:"skill_gaps"`
	Explanation        Explanation         `json:"explanation"`
}

// CompanyResponse is returned by GET /api/v1/companies/readiness.
type CompanyResponse struct {
	Companies     []CompanyEntry `json:"companies"`
	Overall       int            `json:"overall"`
	TopSkills     []SkillSummary `json:"top_skills"`
	WeakestSkills []SkillSummary `json:"weakest_skills"`
	SkillGaps     []SkillGap     `json:"skill_gaps"`
	Explanations  []Explanation  `json:"explanations"`
	Version       string         `json:"version"`
}

// DashboardResponse aggregates skill and company readiness for validation.
type DashboardResponse struct {
	SkillMastery     SkillMasteryResponse `json:"skill_mastery"`
	CompanyReadiness CompanyResponse      `json:"company_readiness"`
}

// CompanyComparison compares V1 and V2 readiness for one company.
type CompanyComparison struct {
	Company string `json:"company"`
	V1Score int    `json:"v1_score"`
	V2Score int    `json:"v2_score"`
	Delta   int    `json:"delta"`
}

// ComparisonSnapshot holds side-by-side V1 and V2 scores.
type ComparisonSnapshot struct {
	OverallV1    int                 `json:"overall_v1"`
	OverallV2    int                 `json:"overall_v2"`
	OverallDelta int                 `json:"overall_delta"`
	ByCompany    []CompanyComparison `json:"by_company"`
	V1Formula    string              `json:"v1_formula"`
	V2Formula    string              `json:"v2_formula"`
}

// ValidationResponse is returned by validation and internal compare APIs.
type ValidationResponse struct {
	ReadinessV2Enabled bool               `json:"readiness_v2_enabled"`
	V1                 V1Snapshot         `json:"v1"`
	V2                 DashboardResponse  `json:"v2"`
	Comparison         ComparisonSnapshot `json:"comparison"`
}
