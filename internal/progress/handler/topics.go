package handler

import (
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// GetTopicMastery handles GET /api/v1/progress/topics.
func (h *ReadinessHandler) GetTopicMastery(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	topics, err := h.readiness.GetTopicMastery(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, topics)
}
