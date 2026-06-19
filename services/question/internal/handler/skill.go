package handler

import (
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/shared/response"
)

// SkillHandler serves skill graph endpoints.
type SkillHandler struct {
	skills *service.SkillService
}

// NewSkillHandler creates a SkillHandler.
func NewSkillHandler(skills *service.SkillService) *SkillHandler {
	return &SkillHandler{skills: skills}
}

// ListSkills handles GET /api/v1/skills.
func (h *SkillHandler) ListSkills(w http.ResponseWriter, r *http.Request) {
	tree, err := h.skills.ListSkillTree(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, tree)
}

// GetSkill handles GET /api/v1/skills/{slug}.
func (h *SkillHandler) GetSkill(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	skill, err := h.skills.GetSkillBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, service.ErrSkillNotFound) {
			response.Error(w, http.StatusNotFound, constants.ErrSkillNotFound, "skill not found")
			return
		}
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, skill)
}

// GetQuestionSkills handles GET /api/v1/questions/{id}/skills.
func (h *SkillHandler) GetQuestionSkills(w http.ResponseWriter, r *http.Request) {
	questionID := r.PathValue("id")
	mappings, err := h.skills.ListQuestionSkills(r.Context(), questionID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, mappings)
}
