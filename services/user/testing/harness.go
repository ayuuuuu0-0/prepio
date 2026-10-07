package testing

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/services/user/internal/dto"
	"github.com/prepio/prepio/services/user/internal/service"
	"github.com/prepio/prepio/services/user/internal/store"
)

// OnboardingRequest and ProfileResponse re-export the DTOs for integration tests.
type (
	OnboardingRequest = dto.OnboardingRequest
	ProfileResponse   = dto.ProfileResponse
)

// ErrInvalidRequest is the error onboarding returns for bad input.
var ErrInvalidRequest = service.ErrInvalidRequest

// NewOnboardingService wires the onboarding service for integration tests.
func NewOnboardingService(pool *pgxpool.Pool) *service.OnboardingService {
	return service.NewOnboardingService(store.NewUserStore(pool), store.NewCharacterStore(pool))
}
