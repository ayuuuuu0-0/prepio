package dto

// SkillCategoryResponse groups skills under a category for the skill tree API.
type SkillCategoryResponse struct {
	ID     string           `json:"id"`
	Slug   string           `json:"slug"`
	Name   string           `json:"name"`
	Skills []SkillResponse  `json:"skills"`
}

// SkillResponse is a skill node in the skill graph.
type SkillResponse struct {
	ID          string              `json:"id"`
	Slug        string              `json:"slug"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Subskills   []SubskillResponse  `json:"subskills,omitempty"`
}

// SubskillResponse is an atomic mastery unit within a skill.
type SubskillResponse struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	Name string `json:"name"`
}

// QuestionSkillResponse maps a question to a skill and subskill with weight.
type QuestionSkillResponse struct {
	SkillSlug    string  `json:"skill_slug"`
	SubskillSlug string  `json:"subskill_slug"`
	Weight       float64 `json:"weight"`
}
