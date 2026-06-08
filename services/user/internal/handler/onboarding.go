package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/user/internal/dto"
	"github.com/prepio/prepio/services/user/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// OnboardingHandler serves onboarding and profile endpoints.
type OnboardingHandler struct {
	onboarding *service.OnboardingService
}

// NewOnboardingHandler creates an OnboardingHandler.
func NewOnboardingHandler(onboarding *service.OnboardingService) *OnboardingHandler {
	return &OnboardingHandler{onboarding: onboarding}
}

// ListCompanions handles GET /api/v1/companions.
func (h *OnboardingHandler) ListCompanions(w http.ResponseWriter, r *http.Request) {
	resp, err := h.onboarding.ListCompanions(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// CompleteOnboarding handles POST /api/v1/users/onboarding.
func (h *OnboardingHandler) CompleteOnboarding(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	var req dto.OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.onboarding.Complete(r.Context(), userID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// GetProfile handles GET /api/v1/users/profile.
func (h *OnboardingHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.onboarding.GetProfile(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	response.Data(w, http.StatusOK, resp)
}
