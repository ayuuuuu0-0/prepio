package dashboard

import (
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// Handler serves dashboard aggregation endpoints.
type Handler struct {
	dashboard *Service
}

// NewHandler creates a Handler.
func NewHandler(dashboard *Service) *Handler {
	return &Handler{dashboard: dashboard}
}

// GetHome handles GET /api/v1/dashboard/home.
func (h *Handler) GetHome(w http.ResponseWriter, r *http.Request) {
	token := middleware.ExtractBearerToken(r)
	if len(token) == 0 {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.dashboard.GetHome(r.Context(), token)
	if err != nil {
		response.Error(w, http.StatusBadGateway, constants.ErrInternal, "failed to load dashboard")
		return
	}

	response.Data(w, http.StatusOK, resp)
}
