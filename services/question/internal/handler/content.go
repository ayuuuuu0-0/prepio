package handler

import (
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/shared/response"
)

// ContentHandler serves journey content binding endpoints.
type ContentHandler struct {
	content   *service.ContentService
	questions *service.QuestionService
}

// NewContentHandler creates a ContentHandler.
func NewContentHandler(content *service.ContentService, questions *service.QuestionService) *ContentHandler {
	return &ContentHandler{content: content, questions: questions}
}

// GetNodeContent handles GET /api/v1/journey/nodes/{id}/content.
func (h *ContentHandler) GetNodeContent(w http.ResponseWriter, r *http.Request) {
	nodeID := r.PathValue("id")
	content, err := h.content.GetNodeContent(r.Context(), nodeID)
	if err != nil {
		if errors.Is(err, service.ErrNodeNotFound) {
			response.Error(w, http.StatusNotFound, constants.ErrJourneyNodeNotFound, "journey node not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, content)
}
