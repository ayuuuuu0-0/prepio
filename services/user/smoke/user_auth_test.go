package smoke_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prepio/prepio/internal/user/handler"
	"github.com/prepio/prepio/internal/user/mailer"
	"github.com/prepio/prepio/internal/user/service"
	"github.com/prepio/prepio/internal/user/store"
	"github.com/prepio/prepio/shared/jwt"
	"github.com/prepio/prepio/shared/middleware"
	"github.com/prepio/prepio/test/factories"
	"github.com/prepio/prepio/test/testdb"
	"github.com/prepio/prepio/test/testredis"
	"github.com/stretchr/testify/require"
)

func TestUserAuthFlow(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	redisClient, _ := testredis.New(t)

	signer, err := jwt.NewSigner("test-secret-key-for-smoke-tests")
	require.NoError(t, err)

	userStore := store.NewUserStore(pool)
	refreshStore := store.NewRefreshTokenStore(pool)
	deviceStore := store.NewUserDeviceStore(pool)
	mockMailer := &mailer.MockMailer{}

	authHandler := handler.NewAuthHandler(service.NewAuthService(userStore, refreshStore, signer, redisClient, mockMailer))
	userHandler := handler.NewUserHandler(service.NewUserService(userStore, deviceStore))

	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/otp/send", authHandler.SendOTP)
		r.Post("/auth/otp/verify", authHandler.VerifyOTP)
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(signer, redisClient))
			r.Post("/auth/logout", authHandler.Logout)
			r.Get("/users/me", userHandler.GetMe)
		})
	})

	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	factory := factories.NewUserFactory()
	registerBody, _ := json.Marshal(map[string]string{
		"email":    factory.Email,
		"username": factory.Username,
		"password": factory.Password,
		"timezone": factory.Timezone,
	})

	registerResp := postJSON(t, server.URL+"/api/v1/auth/register", registerBody, "")
	require.Equal(t, http.StatusCreated, registerResp.StatusCode)

	auth := decodeMap(t, registerResp)
	require.NotEmpty(t, auth["access_token"])
	require.NotEmpty(t, auth["refresh_token"])

	user := auth["user"].(map[string]any)
	require.Equal(t, factory.Email, user["email"])

	accessToken := auth["access_token"].(string)
	refreshToken := auth["refresh_token"].(string)

	profileResp := getJSON(t, server.URL+"/api/v1/users/me", accessToken)
	require.Equal(t, http.StatusOK, profileResp.StatusCode)

	profile := decodeMap(t, profileResp)
	require.Equal(t, user["id"], profile["id"])

	refreshBody, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
	refreshResp := postJSON(t, server.URL+"/api/v1/auth/refresh", refreshBody, "")
	require.Equal(t, http.StatusOK, refreshResp.StatusCode)

	refreshed := decodeMap(t, refreshResp)
	newAccess := refreshed["access_token"].(string)
	require.NotEmpty(t, newAccess)
	require.NotEqual(t, accessToken, newAccess)

	logoutResp := postJSON(t, server.URL+"/api/v1/auth/logout", nil, newAccess)
	require.Equal(t, http.StatusOK, logoutResp.StatusCode)

	afterLogout := getJSON(t, server.URL+"/api/v1/users/me", newAccess)
	require.Equal(t, http.StatusUnauthorized, afterLogout.StatusCode)
}

func TestUserOTPFlow(t *testing.T) {
	pool, _ := testdb.Start(t)
	testdb.Migrate(t, pool)

	redisClient, _ := testredis.New(t)

	signer, err := jwt.NewSigner("test-secret-key-for-smoke-tests")
	require.NoError(t, err)

	userStore := store.NewUserStore(pool)
	refreshStore := store.NewRefreshTokenStore(pool)
	deviceStore := store.NewUserDeviceStore(pool)
	mockMailer := &mailer.MockMailer{}

	authHandler := handler.NewAuthHandler(service.NewAuthService(userStore, refreshStore, signer, redisClient, mockMailer))
	userHandler := handler.NewUserHandler(service.NewUserService(userStore, deviceStore))

	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/otp/send", authHandler.SendOTP)
		r.Post("/auth/otp/verify", authHandler.VerifyOTP)
		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(signer, redisClient))
			r.Get("/users/me", userHandler.GetMe)
		})
	})

	server := httptest.NewServer(r)
	t.Cleanup(server.Close)

	testEmail := "otpuser@example.com"

	// 1. Send OTP
	sendBody, _ := json.Marshal(map[string]string{"email": testEmail})
	sendResp := postJSON(t, server.URL+"/api/v1/auth/otp/send", sendBody, "")
	require.Equal(t, http.StatusOK, sendResp.StatusCode)
	require.Equal(t, testEmail, mockMailer.LastEmail)
	require.Len(t, mockMailer.LastCode, 6)

	otpCode := mockMailer.LastCode

	// 2. Verify with wrong code should fail
	wrongBody, _ := json.Marshal(map[string]string{
		"email": testEmail,
		"code":  "000000",
	})
	wrongResp := postJSON(t, server.URL+"/api/v1/auth/otp/verify", wrongBody, "")
	require.Equal(t, http.StatusBadRequest, wrongResp.StatusCode)

	// 3. Verify with correct code should succeed and create account
	correctBody, _ := json.Marshal(map[string]string{
		"email":    testEmail,
		"code":     otpCode,
		"username": "otpstudent",
	})
	verifyResp := postJSON(t, server.URL+"/api/v1/auth/otp/verify", correctBody, "")
	require.Equal(t, http.StatusOK, verifyResp.StatusCode)

	auth := decodeMap(t, verifyResp)
	require.NotEmpty(t, auth["access_token"])
	require.NotEmpty(t, auth["refresh_token"])

	user := auth["user"].(map[string]any)
	require.Equal(t, testEmail, user["email"])
	require.Equal(t, "otpstudent", user["username"])

	// 4. Access protected endpoint with token
	accessToken := auth["access_token"].(string)
	profileResp := getJSON(t, server.URL+"/api/v1/users/me", accessToken)
	require.Equal(t, http.StatusOK, profileResp.StatusCode)
}

func postJSON(t *testing.T, url string, body []byte, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func getJSON(t *testing.T, url, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	if len(token) > 0 {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func decodeMap(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	defer resp.Body.Close()
	var envelope struct {
		Data map[string]any `json:"data"`
	}
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&envelope))
	return envelope.Data
}
