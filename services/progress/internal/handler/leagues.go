package handler

import (
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// LeagueHandler serves the weekly league.
type LeagueHandler struct {
	leagues *service.LeagueService
}

// NewLeagueHandler creates a LeagueHandler.
func NewLeagueHandler(leagues *service.LeagueService) *LeagueHandler {
	return &LeagueHandler{leagues: leagues}
}

// GetLeague handles GET /api/v1/progress/league.
func (h *LeagueHandler) GetLeague(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	league, err := h.leagues.GetLeague(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, league)
}
