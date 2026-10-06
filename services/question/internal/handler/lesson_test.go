package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prepio/prepio/constants"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/stretchr/testify/require"
)

func TestWriteLessonErrorMapsEveryRuntimeError(t *testing.T) {
	tests := []struct {
		err    error
		status int
		code   string
	}{
		{service.ErrInvalidRequest, http.StatusBadRequest, constants.ErrInvalidRequest},
		{fmt.Errorf("%w: choice 9 is not an option", service.ErrInvalidAnswer), http.StatusBadRequest, constants.ErrInvalidAnswer},
		{service.ErrLessonNotFound, http.StatusNotFound, constants.ErrLessonNotFound},
		{service.ErrLessonLocked, http.StatusForbidden, constants.ErrLessonLocked},
		{service.ErrAttemptNotFound, http.StatusNotFound, constants.ErrAttemptNotFound},
		{service.ErrStepNotFound, http.StatusNotFound, constants.ErrStepNotFound},
		{service.ErrAttemptNotInProgress, http.StatusConflict, constants.ErrAttemptNotInProgress},
		{service.ErrAttemptIncomplete, http.StatusConflict, constants.ErrAttemptIncomplete},
		{service.ErrInvalidTry, http.StatusConflict, constants.ErrInvalidTry},
		{fmt.Errorf("database exploded"), http.StatusInternalServerError, constants.ErrInternal},
	}
	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			rec := httptest.NewRecorder()
			writeLessonError(rec, tt.err)
			require.Equal(t, tt.status, rec.Code)
			var body struct {
				Error struct{ Code, Message string } `json:"error"`
			}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			require.Equal(t, tt.code, body.Error.Code)
			if tt.status == http.StatusInternalServerError {
				require.NotContains(t, body.Error.Message, "exploded", "internal errors must not leak details")
			}
		})
	}
}

func TestLessonHandlersRequireAuthentication(t *testing.T) {
	h := NewLessonHandler(nil)
	for name, fn := range map[string]http.HandlerFunc{
		"path": h.GetPath, "start": h.StartAttempt, "answer": h.SubmitAnswer, "complete": h.Complete,
	} {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			fn(rec, httptest.NewRequest(http.MethodPost, "/", nil))
			require.Equal(t, http.StatusUnauthorized, rec.Code)
		})
	}
}
