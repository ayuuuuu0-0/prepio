package service

import (
	"context"
	"errors"
	"slices"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/user/internal/dto"
	"github.com/prepio/prepio/services/user/internal/store"
)

// OnboardingService handles onboarding and companion selection.
type OnboardingService struct {
	users      *store.UserStore
	characters *store.CharacterStore
}

// NewOnboardingService creates an OnboardingService.
func NewOnboardingService(users *store.UserStore, characters *store.CharacterStore) *OnboardingService {
	return &OnboardingService{users: users, characters: characters}
}

// ListCompanions returns starter companions for onboarding.
func (s *OnboardingService) ListCompanions(ctx context.Context) ([]dto.CharacterResponse, error) {
	characters, err := s.characters.ListStarters(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.CharacterResponse, 0, len(characters))
	for _, c := range characters {
		if !slices.Contains(constants.StarterCompanionIDs, c.ID) {
			continue
		}
		resp = append(resp, dto.CharacterResponse{ID: c.ID, Name: c.Name, Species: c.Species})
	}
	return resp, nil
}

// Complete stores onboarding choices for the authenticated user: 1-3 focus topics, an
// experience level, and a starter companion.
func (s *OnboardingService) Complete(ctx context.Context, userID string, req dto.OnboardingRequest) (*dto.ProfileResponse, error) {
	if len(req.ExperienceLevel) == 0 || len(req.CompanionID) == 0 {
		return nil, ErrInvalidRequest
	}
	if err := validateFocusTopics(req.FocusTopics); err != nil {
		return nil, err
	}
	if !slices.Contains(constants.ExperienceLevels, req.ExperienceLevel) {
		return nil, ErrInvalidRequest
	}
	if !slices.Contains(constants.StarterCompanionIDs, req.CompanionID) {
		return nil, ErrInvalidRequest
	}

	character, err := s.characters.GetByID(ctx, req.CompanionID)
	if err != nil {
		return nil, err
	}
	if character == nil {
		return nil, ErrInvalidRequest
	}

	if err := s.users.UnlockCharacter(ctx, userID, req.CompanionID); err != nil {
		return nil, err
	}

	if err := s.users.SetFocusTopics(ctx, userID, req.FocusTopics); err != nil {
		if errors.Is(err, store.ErrUnknownTopic) {
			return nil, ErrInvalidRequest
		}
		return nil, err
	}

	user, err := s.users.CompleteOnboarding(ctx, userID, req.ExperienceLevel, req.CompanionID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return s.buildProfile(ctx, user)
}

// GetProfile returns the full profile for the authenticated user.
func (s *OnboardingService) GetProfile(ctx context.Context, userID string) (*dto.ProfileResponse, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return s.buildProfile(ctx, user)
}

func (s *OnboardingService) buildProfile(ctx context.Context, user *store.User) (*dto.ProfileResponse, error) {
	focusTopics, err := s.users.ListFocusTopics(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	resp := &dto.ProfileResponse{
		ID:                  user.ID,
		Email:               user.Email,
		Username:            user.Username,
		Timezone:            user.Timezone,
		OnboardingCompleted: user.OnboardingCompleted,
		FocusTopics:         focusTopics,
	}
	if user.ExperienceLevel != nil {
		resp.ExperienceLevel = *user.ExperienceLevel
	}
	if user.ActiveCharID != nil && len(*user.ActiveCharID) > 0 {
		character, err := s.characters.GetByID(ctx, *user.ActiveCharID)
		if err != nil {
			return nil, err
		}
		if character != nil {
			resp.Companion = &dto.CharacterResponse{
				ID:      character.ID,
				Name:    character.Name,
				Species: character.Species,
			}
		}
	}
	return resp, nil
}

// maxFocusTopics is how many topics a learner may focus on.
const maxFocusTopics = 3

// validateFocusTopics requires 1-3 distinct topic slugs. Whether each slug exists is
// checked against the topic catalog when they are stored.
func validateFocusTopics(slugs []string) error {
	if len(slugs) < 1 || len(slugs) > maxFocusTopics {
		return ErrInvalidRequest
	}
	seen := make(map[string]bool, len(slugs))
	for _, s := range slugs {
		if s == "" || seen[s] {
			return ErrInvalidRequest
		}
		seen[s] = true
	}
	return nil
}
