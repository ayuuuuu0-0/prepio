package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/response"
)

// InternalLessonCompleted handles POST /internal/events/lesson-completed (dev sync).
func (h *StreakHandler) InternalLessonCompleted(w http.ResponseWriter, r *http.Request) {
	var event events.LessonCompleted
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
		return
	}
	if err := h.streaks.ProcessLessonCompleted(r.Context(), event); err != nil {
		log.Printf("streak: process lesson completed: %v", err)
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, map[string]bool{"ok": true})
}
