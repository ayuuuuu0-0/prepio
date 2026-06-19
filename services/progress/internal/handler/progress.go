package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/progress/internal/dto"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// ProgressHandler serves progress endpoints.
type ProgressHandler struct {
	progress  *service.ProgressService
	readiness *service.ReadinessService
}

// NewProgressHandler creates a ProgressHandler.
func NewProgressHandler(progress *service.ProgressService, readiness *service.ReadinessService) *ProgressHandler {
	return &ProgressHandler{progress: progress, readiness: readiness}
}

// GetMe handles GET /api/v1/progress/me.
func (h *ProgressHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.progress.GetMe(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// InternalGetGems handles GET /internal/progress/{userID}/gems.
func (h *ProgressHandler) InternalGetGems(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	balance, err := h.progress.GetGems(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, map[string]int{"gem_balance": balance})
}

// InternalDeductGems handles POST /internal/progress/{userID}/gems/deduct.
func (h *ProgressHandler) InternalDeductGems(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	var req dto.DeductGemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	balance, err := h.progress.DeductGems(r.Context(), userID, req.Amount, req.Reason)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientGems) {
			response.Error(w, http.StatusConflict, constants.ErrInsufficientGems, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, dto.DeductGemsResponse{GemBalance: balance})
}
