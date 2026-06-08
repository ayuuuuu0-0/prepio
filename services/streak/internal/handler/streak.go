package handler

import (
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/streak/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// StreakHandler serves streak endpoints.
type StreakHandler struct {
	streaks *service.StreakService
}

// NewStreakHandler creates a StreakHandler.
func NewStreakHandler(streaks *service.StreakService) *StreakHandler {
	return &StreakHandler{streaks: streaks}
}

// GetMe handles GET /api/v1/streaks/me.
func (h *StreakHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.streaks.GetMe(r.Context(), userID, r.URL.Query().Get("timezone"))
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// PurchaseFreeze handles POST /api/v1/streaks/me/freeze/purchase.
func (h *StreakHandler) PurchaseFreeze(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.streaks.PurchaseFreeze(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	response.Data(w, http.StatusOK, resp)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrFreezeMaxHeld):
		response.Error(w, http.StatusConflict, constants.ErrStreakFreezeMaxHeld, err.Error())
	case errors.Is(err, service.ErrInsufficientGems):
		response.Error(w, http.StatusConflict, constants.ErrStreakFreezeInsufficientGems, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
	}
}
