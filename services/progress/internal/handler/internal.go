package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/response"
)

// InternalQuestionAnswered handles POST /internal/events/question-answered (dev sync).
func (h *ProgressHandler) InternalQuestionAnswered(w http.ResponseWriter, r *http.Request) {
	var event events.QuestionAnswered
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
		return
	}
	if err := h.progress.ProcessQuestionAnswered(r.Context(), event); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	if h.readiness != nil {
		if err := h.readiness.ProcessQuestionAnswered(r.Context(), event); err != nil {
			response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
			return
		}
	}
	response.Data(w, http.StatusOK, map[string]bool{"ok": true})
}

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
