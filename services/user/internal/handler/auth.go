package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/user/internal/dto"
	"github.com/prepio/prepio/services/user/internal/service"
	"github.com/prepio/prepio/shared/auth"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// AuthHandler serves auth endpoints.
type AuthHandler struct {
	auth *service.AuthService
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(auth *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

// Register handles POST /api/v1/auth/register.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Register(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	auth.SetRefreshTokenCookie(w, resp.RefreshToken)
	response.Data(w, http.StatusCreated, resp)
}

// Login handles POST /api/v1/auth/login.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.auth.Login(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	auth.SetRefreshTokenCookie(w, resp.RefreshToken)
	response.Data(w, http.StatusOK, resp)
}

// Refresh handles POST /api/v1/auth/refresh.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	if len(req.RefreshToken) == 0 {
		req.RefreshToken = auth.RefreshTokenFromRequest(r)
	}

	resp, err := h.auth.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	auth.SetRefreshTokenCookie(w, resp.RefreshToken)
	response.Data(w, http.StatusOK, resp)
}

// Logout handles POST /api/v1/auth/logout.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	token := middleware.ExtractBearerToken(r)
	if len(token) == 0 {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	if err := h.auth.Logout(r.Context(), token); err != nil {
		writeServiceError(w, err)
		return
	}

	auth.ClearRefreshTokenCookie(w)
	response.Data(w, http.StatusOK, map[string]bool{"logged_out": true})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidRequest):
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(w, http.StatusUnauthorized, constants.ErrInvalidCredentials, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		response.Error(w, http.StatusConflict, constants.ErrEmailTaken, err.Error())
	case errors.Is(err, service.ErrUsernameTaken):
		response.Error(w, http.StatusConflict, constants.ErrUsernameTaken, err.Error())
	case errors.Is(err, service.ErrInvalidToken):
		response.Error(w, http.StatusUnauthorized, constants.ErrInvalidToken, err.Error())
	case errors.Is(err, service.ErrRefreshTokenInvalid):
		response.Error(w, http.StatusUnauthorized, constants.ErrRefreshTokenInvalid, err.Error())
	case errors.Is(err, service.ErrUserNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrUserNotFound, err.Error())
	case errors.Is(err, service.ErrDeviceNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrDeviceNotFound, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
	}
}
