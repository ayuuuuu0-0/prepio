package smoke_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prepio/prepio/services/question/internal/handler"
	"github.com/prepio/prepio/services/question/internal/service"
	"github.com/prepio/prepio/services/question/internal/store"
	"github.com/prepio/prepio/shared/events"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/test/fakes"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestDailyPaperAndSubmitEmitsEvent(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	redisClient, _ := testredis.New(t)
	publisher := &fakes.KafkaProducer{}

	userID := seedUser(t, pool)
	seedQuestion(t, pool)

	signer, err := jwt.NewSigner("question-smoke-secret")
	require.NoError(t, err)

	questionService := service.NewQuestionService(
		store.NewQuestionStore(pool),
		store.NewDailyPaperStore(pool),
		store.NewHistoryStore(pool),
		store.NewUserStore(pool),
		redisClient,
		publisher,
	)
	questionHandler := handler.NewQuestionHandler(questionService)

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.Auth(signer, redisClient))
		r.Get("/questions/daily", questionHandler.GetDaily)
		r.Post("/questions/{id}/submit", questionHandler.Submit)
	})
	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	token, _, _, err := signer.SignAccessToken(userID)
	require.NoError(t, err)

	dailyResp := getAuth(t, server.URL+"/api/v1/questions/daily", token)
	require.Equal(t, http.StatusOK, dailyResp.StatusCode)
	daily := decode(t, dailyResp)
	require.NotEmpty(t, daily["session_id"])

	questions := daily["questions"].([]any)
	require.NotEmpty(t, questions)
	first := questions[0].(map[string]any)
	questionID := first["id"].(string)

	submitBody, _ := json.Marshal(map[string]any{
		"session_id":         daily["session_id"],
		"answer":             "use hash map approach with O(n) time and O(n) space complexity",
		"time_spent_seconds": 120,
	})
	submitResp := postAuth(t, fmt.Sprintf("%s/api/v1/questions/%s/submit", server.URL, questionID), submitBody, token)
	require.Equal(t, http.StatusOK, submitResp.StatusCode)
	submit := decode(t, submitResp)
	require.Equal(t, true, submit["correct"])

	last := publisher.Last()
	require.NotNil(t, last)
	require.Equal(t, events.TopicQuestionAnswered, last.Topic)

	var event events.QuestionAnswered
	require.NoError(t, json.Unmarshal(last.Payload, &event))
	require.Equal(t, userID, event.UserID)
	require.Equal(t, questionID, event.QuestionID)
	require.True(t, event.Correct)
}

func seedUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	ctx := context.Background()
	var userID string
	err := pool.QueryRow(ctx, `
		INSERT INTO users (email, username, password_hash)
		VALUES ('q@test.com', 'quser', 'hash')
		RETURNING id`).Scan(&userID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO character_unlocks (user_id, character_id)
		VALUES ($1, 'a0000000-0000-4000-8000-000000000001')`, userID)
	require.NoError(t, err)
	return userID
}

func seedQuestion(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	ctx := context.Background()
	var questionID string
	err := pool.QueryRow(ctx, `
		INSERT INTO questions (body, round_type, difficulty, answer_guide, status, source)
		VALUES (
			'Given an array, return two sum indices',
			'dsa', 'easy',
			'hash map approach with O(n) time and O(n) space',
			'approved', 'manual'
		) RETURNING id`).Scan(&questionID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `INSERT INTO question_tags (question_id, company) VALUES ($1, 'google')`, questionID)
	require.NoError(t, err)
	return questionID
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

func postAuth(t *testing.T, url string, body []byte, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decode(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&envelope))
	return envelope.Data
}
