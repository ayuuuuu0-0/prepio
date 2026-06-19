package dto

// NodeSkillResponse describes a skill taught or evaluated by a journey node.
type NodeSkillResponse struct {
	SkillSlug string `json:"skill_slug"`
	SkillName string `json:"skill_name"`
	IsPrimary bool   `json:"is_primary"`
}

// NodePoolResponse describes a question pool bound to a journey node.
type NodePoolResponse struct {
	PoolSlug          string `json:"pool_slug"`
	PoolName          string `json:"pool_name"`
	SkillSlug         string `json:"skill_slug"`
	SelectionStrategy string `json:"selection_strategy"`
	QuestionsRequired int    `json:"questions_required"`
	QuestionCount     int    `json:"question_count"`
}

// NodeContentResponse is returned by GET /api/v1/journey/nodes/{id}/content.
type NodeContentResponse struct {
	NodeID string              `json:"node_id"`
	Slug   string              `json:"slug,omitempty"`
	Label  string              `json:"label"`
	Skills []NodeSkillResponse `json:"skills"`
	Pools  []NodePoolResponse  `json:"pools"`
}
