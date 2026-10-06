package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// LessonHandler serves lesson reward endpoints.
type LessonHandler struct {
	lessons *service.LessonService
}

// NewLessonHandler creates a LessonHandler.
func NewLessonHandler(lessons *service.LessonService) *LessonHandler {
	return &LessonHandler{lessons: lessons}
}

// GetAttemptRewards handles GET /api/v1/progress/attempts/{attemptID}/rewards.
func (h *LessonHandler) GetAttemptRewards(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	resp, found, err := h.lessons.GetAttemptRewards(r.Context(), userID, r.PathValue("attemptID"))
	if errors.Is(err, service.ErrInvalidRequest) {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid attempt id")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	if !found {
		response.Error(w, http.StatusNotFound, constants.ErrProgressNotFound, "rewards not available yet")
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// InternalLessonCompleted handles POST /internal/events/lesson-completed (dev sync).
func (h *LessonHandler) InternalLessonCompleted(w http.ResponseWriter, r *http.Request) {
	var event events.LessonCompleted
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
		return
	}
	if err := h.lessons.ProcessLessonCompleted(r.Context(), event); err != nil {
		if errors.Is(err, service.ErrInvalidRequest) {
			response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid event")
			return
		}
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, map[string]bool{"ok": true})
}
