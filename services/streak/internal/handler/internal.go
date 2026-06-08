package handler

import (
	"encoding/json"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/response"
)

// InternalQuestionAnswered handles POST /internal/events/question-answered (dev sync).
func (h *StreakHandler) InternalQuestionAnswered(w http.ResponseWriter, r *http.Request) {
	var event events.QuestionAnswered
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid body")
		return
	}
	if err := h.streaks.ProcessQuestionAnswered(r.Context(), event); err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, map[string]bool{"ok": true})
}
