package handler

import (
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// ReadinessHandler serves readiness V2 endpoints.
type ReadinessHandler struct {
	readiness *service.ReadinessService
}

// NewReadinessHandler creates a ReadinessHandler.
func NewReadinessHandler(readiness *service.ReadinessService) *ReadinessHandler {
	return &ReadinessHandler{readiness: readiness}
}

// GetSkillReadiness handles GET /api/v1/skills/readiness.
func (h *ReadinessHandler) GetSkillReadiness(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.readiness.GetSkillReadiness(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// GetCompanyReadiness handles GET /api/v1/companies/readiness.
func (h *ReadinessHandler) GetCompanyReadiness(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.readiness.GetCompanyReadiness(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// GetReadinessDashboard handles GET /api/v1/readiness/dashboard.
func (h *ReadinessHandler) GetReadinessDashboard(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.readiness.GetReadinessDashboard(r.Context(), userID)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, err.Error())
			return
		}
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, resp)
}
