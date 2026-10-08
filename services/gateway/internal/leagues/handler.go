package leagues

import (
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// Handler serves the league endpoint.
type Handler struct {
	leagues *Service
}

// NewHandler creates a Handler.
func NewHandler(leagues *Service) *Handler {
	return &Handler{leagues: leagues}
}

// GetLeague handles GET /api/v1/league.
func (h *Handler) GetLeague(w http.ResponseWriter, r *http.Request) {
	token := middleware.ExtractBearerToken(r)
	if len(token) == 0 {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	league, err := h.leagues.GetLeague(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusBadGateway, constants.ErrInternal, "failed to load league")
		return
	}
	response.Data(w, http.StatusOK, league)
}
