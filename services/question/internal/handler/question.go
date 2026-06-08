package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// QuestionHandler serves question endpoints.
type QuestionHandler struct {
	questions *service.QuestionService
}

// NewQuestionHandler creates a QuestionHandler.
func NewQuestionHandler(questions *service.QuestionService) *QuestionHandler {
	return &QuestionHandler{questions: questions}
}

// GetDaily handles GET /api/v1/questions/daily.
func (h *QuestionHandler) GetDaily(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	timezone := r.URL.Query().Get("timezone")
	resp, err := h.questions.GetDailyPaper(r.Context(), userID, timezone)
	if err != nil {
		writeError(w, err)
		return
	}

	response.Data(w, http.StatusOK, resp)
}

// Submit handles POST /api/v1/questions/{id}/submit.
func (h *QuestionHandler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	questionID := r.PathValue("id")
	var req dto.SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.questions.SubmitAnswer(r.Context(), userID, questionID, req)
	if err != nil {
		writeError(w, err)
		return
	}

	response.Data(w, http.StatusOK, resp)
}

// ListCompanies handles GET /api/v1/questions/companies.
func (h *QuestionHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := h.questions.ListCompanies(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
		return
	}
	response.Data(w, http.StatusOK, companies)
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidRequest):
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, err.Error())
	case errors.Is(err, service.ErrQuestionNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrQuestionNotFound, err.Error())
	case errors.Is(err, service.ErrQuestionNotInSession):
		response.Error(w, http.StatusBadRequest, constants.ErrQuestionNotInSession, err.Error())
	case errors.Is(err, service.ErrAnswerAlreadySubmitted):
		response.Error(w, http.StatusConflict, constants.ErrAnswerAlreadySubmitted, err.Error())
	default:
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
	}
}
