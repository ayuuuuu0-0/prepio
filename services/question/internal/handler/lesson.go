package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/dto"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
)

// maxAnswerBody bounds a step-answer request body.
const maxAnswerBody = 64 << 10

// LessonHandler serves the lesson runtime endpoints.
type LessonHandler struct {
	lessons *service.LessonService
}

// NewLessonHandler creates a LessonHandler.
func NewLessonHandler(lessons *service.LessonService) *LessonHandler {
	return &LessonHandler{lessons: lessons}
}

// GetPath handles GET /api/v1/path.
func (h *LessonHandler) GetPath(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	resp, err := h.lessons.GetPath(r.Context(), userID)
	if err != nil {
		writeLessonError(w, err)
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// StartAttempt handles POST /api/v1/lessons/{id}/attempts.
func (h *LessonHandler) StartAttempt(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	resp, err := h.lessons.StartAttempt(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		writeLessonError(w, err)
		return
	}
	status := http.StatusCreated
	if resp.Resumed {
		status = http.StatusOK
	}
	response.Data(w, status, resp)
}

// SubmitAnswer handles POST /api/v1/attempts/{id}/steps/{stepId}/answer.
func (h *LessonHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}

	var req dto.AnswerRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxAnswerBody))
	if err := dec.Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, "invalid request body")
		return
	}

	resp, err := h.lessons.SubmitAnswer(r.Context(), userID, r.PathValue("id"), r.PathValue("stepId"), req)
	if err != nil {
		writeLessonError(w, err)
		return
	}
	response.Data(w, http.StatusOK, resp)
}

// Complete handles POST /api/v1/attempts/{id}/complete.
func (h *LessonHandler) Complete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, constants.ErrUnauthorized, "authorization required")
		return
	}
	resp, err := h.lessons.Complete(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		writeLessonError(w, err)
		return
	}
	response.Data(w, http.StatusOK, resp)
}

func writeLessonError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidRequest):
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidRequest, err.Error())
	case errors.Is(err, service.ErrInvalidAnswer):
		response.Error(w, http.StatusBadRequest, constants.ErrInvalidAnswer, err.Error())
	case errors.Is(err, service.ErrLessonNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrLessonNotFound, err.Error())
	case errors.Is(err, service.ErrLessonLocked):
		response.Error(w, http.StatusForbidden, constants.ErrLessonLocked, err.Error())
	case errors.Is(err, service.ErrAttemptNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrAttemptNotFound, err.Error())
	case errors.Is(err, service.ErrStepNotFound):
		response.Error(w, http.StatusNotFound, constants.ErrStepNotFound, err.Error())
	case errors.Is(err, service.ErrAttemptNotInProgress):
		response.Error(w, http.StatusConflict, constants.ErrAttemptNotInProgress, err.Error())
	case errors.Is(err, service.ErrAttemptIncomplete):
		response.Error(w, http.StatusConflict, constants.ErrAttemptIncomplete, err.Error())
	case errors.Is(err, service.ErrInvalidTry):
		response.Error(w, http.StatusConflict, constants.ErrInvalidTry, err.Error())
	default:
		log.Printf("lesson: unhandled error: %v", err)
		response.Error(w, http.StatusInternalServerError, constants.ErrInternal, "internal error")
	}
}
