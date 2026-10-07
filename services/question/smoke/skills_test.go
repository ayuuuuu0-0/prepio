package smoke_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prepio/prepio/services/question/internal/handler"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

func TestListSkillsEndpoint(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	redisClient, _ := testredis.New(t)
	userID := seedUser(t, pool)

	signer, err := jwt.NewSigner("skills-smoke-secret")
	require.NoError(t, err)

	skillService := service.NewSkillService(store.NewSkillStore(pool))
	skillHandler := handler.NewSkillHandler(skillService)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(signer, redisClient))
		r.Get("/skills", skillHandler.ListSkills)
		r.Get("/skills/{slug}", skillHandler.GetSkill)
	})
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	token, _, _, err := signer.SignAccessToken(userID)
	require.NoError(t, err)

	listResp := getAuth(t, server.URL+"/api/v1/skills", token)
	require.Equal(t, http.StatusOK, listResp.StatusCode)
	list := decodeDataArray(t, listResp)
	require.NotEmpty(t, list)

	skillResp := getAuth(t, server.URL+"/api/v1/skills/arrays", token)
	require.Equal(t, http.StatusOK, skillResp.StatusCode)
	skill := decodeDataMap(t, skillResp)
	require.Equal(t, "arrays", skill["slug"])
}

func decodeDataArray(t *testing.T, resp *http.Response) []any {
	t.Helper()
	defer resp.Body.Close()
	var envelope struct {
		Data []any `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&envelope))
	return envelope.Data
}

func decodeDataMap(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&envelope))
	return envelope.Data
}

func seedUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	var userID string
	err := pool.QueryRow(context.Background(), `
		INSERT INTO users (email, username, password_hash)
		VALUES ('skills@test.com', 'skillsuser', 'hash')
		RETURNING id`).Scan(&userID)
	require.NoError(t, err)
	return userID
}

func getAuth(t *testing.T, url, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
