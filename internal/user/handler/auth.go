package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/internal/user/dto"
	"github.com/prepio/prepio/internal/user/service"
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

// SendOTP handles POST /api/v1/auth/otp/send.
func (h *AuthHandler) SendOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.SendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.auth.SendOTP(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	response.Data(w, http.StatusOK, resp)
}

// VerifyOTP handles POST /api/v1/auth/otp/verify.
func (h *AuthHandler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	var req dto.VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.auth.VerifyOTP(r.Context(), req)
	if err != nil {
		writeServiceError(w, err)
		return
	}

	auth.SetRefreshTokenCookie(w, r, resp.RefreshToken)
	response.Data(w, http.StatusOK, resp)
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

	auth.SetRefreshTokenCookie(w, r, resp.RefreshToken)
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

	auth.SetRefreshTokenCookie(w, r, resp.RefreshToken)
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

	auth.SetRefreshTokenCookie(w, r, resp.RefreshToken)
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

	auth.ClearRefreshTokenCookie(w, r)
	response.Data(w, http.StatusOK, map[string]bool{"logged_out": true})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidRequest):
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		response.Error(w, http.StatusUnauthorized, constants.ErrInvalidCredentials, "Incorrect password. Please try again.")
	case errors.Is(err, service.ErrInvalidOTP):
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidOTP, err.Error())
	case errors.Is(err, service.ErrOTPCooldown):
		response.Error(w, http.StatusTooManyRequests, constants.ErrOTPCooldown, err.Error())
	case errors.Is(err, service.ErrTooManyAttempts):
		response.Error(w, http.StatusTooManyRequests, constants.ErrTooManyAttempts, err.Error())
	case errors.Is(err, service.ErrEmailTaken):
		response.Error(w, http.StatusConflict, constants.ErrEmailTaken, "An account with this email already exists.")
	case errors.Is(err, service.ErrUsernameTaken):
		response.Error(w, http.StatusConflict, constants.ErrUsernameTaken, "This username is already taken.")
	case errors.Is(err, service.ErrInvalidToken):
		response.Error(w, http.StatusUnauthorized, constants.ErrInvalidToken, err.Error())
	case errors.Is(err, service.ErrRefreshTokenInvalid):
		response.Error(w, http.StatusUnauthorized, constants.ErrRefreshTokenInvalid, err.Error())
	case errors.Is(err, service.ErrUserNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrUserNotFound, "Account does not exist. Please sign up.")
	case errors.Is(err, service.ErrDeviceNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrDeviceNotFound, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
	}
}
