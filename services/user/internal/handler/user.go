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

// UserHandler serves user profile and device endpoints.
type UserHandler struct {
	users *service.UserService
}

// NewUserHandler creates a UserHandler.
func NewUserHandler(users *service.UserService) *UserHandler {
	return &UserHandler{users: users}
}

// GetMe handles GET /api/v1/users/me.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	resp, err := h.users.GetProfile(r.Context(), userID)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.Data(w, http.StatusOK, resp)
}

// UpdateMe handles PATCH /api/v1/users/me.
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	var req dto.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.users.UpdateProfile(r.Context(), userID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.Data(w, http.StatusOK, resp)
}

// RegisterDevice handles POST /api/v1/users/me/devices.
func (h *UserHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	var req dto.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.users.RegisterDevice(r.Context(), userID, req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.Data(w, http.StatusCreated, resp)
}

// DeleteDevice handles DELETE /api/v1/users/me/devices/{deviceID}.
func (h *UserHandler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	deviceID := r.PathValue("deviceID")
	if err := h.users.DeleteDevice(r.Context(), userID, deviceID); err != nil {
		writeServiceError(w, err)
		return
	}

	response.Data(w, http.StatusOK, map[string]bool{"deleted": true})
}
