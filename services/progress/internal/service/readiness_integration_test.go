package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prepio/prepio/services/progress/internal/handler"
	"github.com/prepio/prepio/services/progress/internal/service"
	"github.com/prepio/prepio/services/progress/internal/store"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

func TestReadinessV2Endpoints(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	ctx := context.Background()
	readinessStore := store.NewReadinessStore(pool)
	readinessService := service.NewReadinessService(readinessStore)

	var userID string
	require.NoError(t, pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('rv2@test.com', 'rv2user', 'hash') RETURNING id`).Scan(&userID))

	_, err := pool.Exec(ctx, `
		INSERT INTO user_skill_scores (user_id, skill_id, mastery, attempts, last_practiced_at)
		VALUES ($1, 'b2000001-0000-4000-8000-000000000017', 42, 2, now())`, userID)
	require.NoError(t, err)

	redisClient, _ := testredis.New(t)
	signer, err := jwt.NewSigner("readiness-v2-smoke")
	require.NoError(t, err)

	readinessHandler := handler.NewReadinessHandler(readinessService)
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(signer, redisClient))
		r.Get("/skills/readiness", readinessHandler.GetSkillReadiness)
	})
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	token, _, _, err := signer.SignAccessToken(userID)
	require.NoError(t, err)

	skillResp := getReadinessAuth(t, server.URL+"/api/v1/skills/readiness", token)
	require.Equal(t, http.StatusOK, skillResp.StatusCode)
	skills := decodeReadinessEnvelope(t, skillResp)
	require.Equal(t, "v2", skills["version"])
	require.NotEmpty(t, skills["skills"])
}

func getReadinessAuth(t *testing.T, url, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decodeReadinessEnvelope(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&envelope))
	return envelope.Data
}
