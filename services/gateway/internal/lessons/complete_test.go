package lessons_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/prepio/prepio/services/gateway/internal/lessons"
	"github.com/stretchr/testify/require"
)

func serve(t *testing.T, question, progress http.HandlerFunc) *httptest.Server {
	t.Helper()
	q := httptest.NewServer(question)
	p := httptest.NewServer(progress)
	t.Cleanup(q.Close)
	t.Cleanup(p.Close)

	h := lessons.NewHandler(q.URL, p.URL)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/attempts/{id}/complete", h.Complete)
	gw := httptest.NewServer(mux)
	t.Cleanup(gw.Close)
	return gw
}

func complete(t *testing.T, gw *httptest.Server, token string) (*http.Response, map[string]any) {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, gw.URL+"/api/v1/attempts/abc/complete", nil)
	require.NoError(t, err)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	t.Cleanup(func() { _ = res.Body.Close() })
	var body map[string]any
	_ = json.NewDecoder(res.Body).Decode(&body)
	return res, body
}

func TestCompleteAttachesRewards(t *testing.T) {
	var gotAuth string
	gw := serve(t,
		func(w http.ResponseWriter, r *http.Request) {
			gotAuth = r.Header.Get("Authorization")
			_, _ = w.Write([]byte(`{"data":{"attempt_id":"abc","accuracy":1}}`))
		},
		func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/progress/attempts/abc/rewards", r.URL.Path)
			_, _ = w.Write([]byte(`{"data":{"xp_awarded":30}}`))
		})

	res, body := complete(t, gw, "tok")
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, "Bearer tok", gotAuth, "the user's token is forwarded")

	data := body["data"].(map[string]any)
	require.Equal(t, "abc", data["attempt_id"])
	require.Equal(t, false, data["rewards_pending"])
	require.EqualValues(t, 30, data["rewards"].(map[string]any)["xp_awarded"])
}

func TestCompleteRetriesUntilProgressHasProcessedTheEvent(t *testing.T) {
	var calls atomic.Int32
	gw := serve(t,
		func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{"data":{"attempt_id":"abc"}}`)) },
		func(w http.ResponseWriter, _ *http.Request) {
			if calls.Add(1) < 3 {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"error":{"code":"progress_not_found","message":"x"}}`))
				return
			}
			_, _ = w.Write([]byte(`{"data":{"xp_awarded":5}}`))
		})

	_, body := complete(t, gw, "tok")
	data := body["data"].(map[string]any)
	require.Equal(t, false, data["rewards_pending"])
	require.EqualValues(t, 3, calls.Load())
}

func TestCompleteReportsPendingRewards(t *testing.T) {
	gw := serve(t,
		func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{"data":{"attempt_id":"abc"}}`)) },
		func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNotFound) })

	res, body := complete(t, gw, "tok")
	require.Equal(t, http.StatusOK, res.StatusCode, "completion succeeds even if rewards are still pending")
	data := body["data"].(map[string]any)
	require.Equal(t, true, data["rewards_pending"])
	require.Nil(t, data["rewards"])
}

func TestCompletePassesThroughLessonErrors(t *testing.T) {
	progressCalled := false
	gw := serve(t,
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"error":{"code":"attempt_incomplete","message":"not done"}}`))
		},
		func(w http.ResponseWriter, _ *http.Request) { progressCalled = true })

	res, body := complete(t, gw, "tok")
	require.Equal(t, http.StatusConflict, res.StatusCode)
	require.Equal(t, "attempt_incomplete", body["error"].(map[string]any)["code"])
	require.False(t, progressCalled)
}

func TestCompleteRequiresAuthorization(t *testing.T) {
	gw := serve(t, func(http.ResponseWriter, *http.Request) {}, func(http.ResponseWriter, *http.Request) {})
	res, _ := complete(t, gw, "")
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}
