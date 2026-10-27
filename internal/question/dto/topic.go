package dto

// TopicResponse is a topic in the catalog (GET /api/v1/topics).
type TopicResponse struct {
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
