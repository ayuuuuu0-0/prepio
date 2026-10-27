package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/response"
)

// InternalStreakUpdated handles POST /internal/events/streak-updated (dev sync).
func (h *ProgressHandler) InternalStreakUpdated(w http.ResponseWriter, r *http.Request) {
	var event events.StreakUpdated
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
		return
	}
	if err := h.progress.ProcessStreakUpdated(r.Context(), event); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, map[string]bool{"ok": true})
}
