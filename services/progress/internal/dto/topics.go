package dto

// TopicSkillResponse is one skill's mastery inside a topic. Mastery is nil until the
// skill has been practiced ("not started", never shown as a failure).
type TopicSkillResponse struct {
	Slug    string `json:"slug"`
	Name    string `json:"name"`
	Mastery *int   `json:"mastery"`
}

// TopicMasteryResponse is a user's readiness in one topic: the mean mastery of the skills
// they have started, and how much of the topic they have started.
type TopicMasteryResponse struct {
	Slug          string               `json:"slug"`
	Name          string               `json:"name"`
	Description   string               `json:"description"`
	Mastery       *int                 `json:"mastery"`
	SkillsStarted int                  `json:"skills_started"`
	SkillsTotal   int                  `json:"skills_total"`
	Skills        []TopicSkillResponse `json:"skills"`
}
