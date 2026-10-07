package dto

// OnboardingRequest is the body for POST /api/v1/users/onboarding.
type OnboardingRequest struct {
	FocusTopics     []string `json:"focus_topics"`
	ExperienceLevel string `json:"experience_level"`
	CompanionID     string `json:"companion_id"`
}

// CharacterResponse is the public companion shape.
type CharacterResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Species string `json:"species"`
}

// ProfileResponse is the full user profile returned by GET /api/v1/users/profile.
type ProfileResponse struct {
	ID                  string             `json:"id"`
	Email               string             `json:"email"`
	Username            string             `json:"username"`
	Timezone            string             `json:"timezone,omitempty"`
	ExperienceLevel     string             `json:"experience_level,omitempty"`
	OnboardingCompleted bool               `json:"onboarding_completed"`
	FocusTopics         []string           `json:"focus_topics"`
	Companion           *CharacterResponse `json:"companion,omitempty"`
}
