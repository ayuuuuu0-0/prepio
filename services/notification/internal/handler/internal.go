package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/notification/internal/service"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/response"
)

// NotificationHandler serves internal dev event endpoints.
type NotificationHandler struct {
	notifications *service.NotificationService
}

// NewNotificationHandler creates a NotificationHandler.
func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notifications: svc}
}

// InternalProgressUpdated handles POST /internal/events/progress-updated (dev sync).
func (h *NotificationHandler) InternalProgressUpdated(w http.ResponseWriter, r *http.Request) {
	var event events.ProgressUpdated
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
		return
	}
	if err := h.notifications.HandleProgressUpdated(r.Context(), event); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, map[string]bool{"ok": true})
}

// InternalStreakUpdated handles POST /internal/events/streak-updated (dev sync).
func (h *NotificationHandler) InternalStreakUpdated(w http.ResponseWriter, r *http.Request) {
	var event events.StreakUpdated
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
		return
	}
	if err := h.notifications.HandleStreakUpdated(r.Context(), event); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, map[string]bool{"ok": true})
}
