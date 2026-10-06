package handler

import (
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/response"
)

// ListTopics handles GET /api/v1/topics.
func (h *SkillHandler) ListTopics(w http.ResponseWriter, r *http.Request) {
	topics, err := h.skills.ListTopics(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, topics)
}
