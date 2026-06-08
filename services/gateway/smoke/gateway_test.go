package smoke_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prepio/prepio/shared/proxy"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/shared/response"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

func TestGatewayProxiesAndEnforcesRateLimit(t *testing.T) {
	redisClient, _ := testredis.New(t)

	signer, err := jwt.NewSigner("gateway-smoke-test-secret")
	require.NoError(t, err)

	var upstreamCalls int
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upstreamCalls++
		switch r.URL.Path {
		case "/api/v1/auth/register":
			response.Data(w, http.StatusCreated, map[string]string{"status": "registered"})
		case "/api/v1/users/me":
			response.Data(w, http.StatusOK, map[string]string{"id": "user-1"})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(upstream.Close)

	upstreamProxy, err := proxy.New(upstream.URL)
	require.NoError(t, err)

	gateway := chi.NewRouter()
	fixedKey := func(*http.Request) string { return "test-client" }

	gateway.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(middleware.RateLimit(redisClient, 3, fixedKey))
			r.Post("/auth/register", upstreamProxy.ServeHTTP)
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(signer, redisClient))
			r.Get("/users/me", upstreamProxy.ServeHTTP)
		})
	})
	gatewayServer := httptest.NewServer(gateway)
	t.Cleanup(gatewayServer.Close)

	body := []byte(`{"email":"a@test.com","username":"a","password":"password123"}`)

	resp := post(t, gatewayServer.URL+"/api/v1/auth/register", body)
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()
	require.Equal(t, 1, upstreamCalls)

	accessToken, _, _, err := signer.SignAccessToken("user-1")
	require.NoError(t, err)

	profile := get(t, gatewayServer.URL+"/api/v1/users/me", accessToken)
	require.Equal(t, http.StatusOK, profile.StatusCode)
	profile.Body.Close()
	require.Equal(t, 2, upstreamCalls)

	for i := 0; i < 2; i++ {
		r := post(t, gatewayServer.URL+"/api/v1/auth/register", body)
		require.Equal(t, http.StatusCreated, r.StatusCode)
		r.Body.Close()
	}
	blocked := post(t, gatewayServer.URL+"/api/v1/auth/register", body)
	require.Equal(t, http.StatusTooManyRequests, blocked.StatusCode)
	blocked.Body.Close()
	require.Equal(t, 4, upstreamCalls)
}

func post(t *testing.T, url string, body []byte) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func get(t *testing.T, url, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}
