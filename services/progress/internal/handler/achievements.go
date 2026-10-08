package handler

import (
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// AchievementHandler serves achievements.
type AchievementHandler struct {
	achievements *service.AchievementService
}

// NewAchievementHandler creates an AchievementHandler.
func NewAchievementHandler(achievements *service.AchievementService) *AchievementHandler {
	return &AchievementHandler{achievements: achievements}
}

// List handles GET /api/v1/progress/achievements.
func (h *AchievementHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	list, err := h.achievements.List(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, list)
}
