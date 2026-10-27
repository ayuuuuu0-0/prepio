package leagues

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const progressLeague = `{"data":{"week_start":"2026-10-05","ends_at":"2026-10-12T00:00:00Z","joined":true,
	"tier":{"index":1,"slug":"silver","name":"Silver"},
	"tiers":[{"index":0,"slug":"bronze","name":"Bronze"},{"index":1,"slug":"silver","name":"Silver"}],
	"my_rank":2,"standings":[
		{"rank":1,"user_id":"u1","weekly_xp":40,"zone":"promote","is_me":false},
		{"rank":2,"user_id":"u2","weekly_xp":25,"zone":"promote","is_me":true}],
	"last_result":null}}`

func upstreams(t *testing.T, cards string, cardStatus int) (progress, user *httptest.Server, gotIDs *[]string) {
	t.Helper()
	ids := []string{}
	progress = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/progress/league", r.URL.Path)
		require.Equal(t, "Bearer tok", r.Header.Get("Authorization"))
		_, _ = w.Write([]byte(progressLeague))
	}))
	user = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/users/public-cards", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "Bearer tok", r.Header.Get("Authorization"))
		var body struct {
			UserIDs []string `json:"user_ids"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		ids = body.UserIDs
		w.WriteHeader(cardStatus)
		_, _ = w.Write([]byte(cards))
	}))
	t.Cleanup(progress.Close)
	t.Cleanup(user.Close)
	return progress, user, &ids
}

func TestGetLeagueAttachesPublicCards(t *testing.T) {
	progress, user, ids := upstreams(t, `{"data":[
		{"id":"u2","username":"bo","companion_name":"Pip","companion_species":"red_panda"},
		{"id":"u1","username":"al","companion_name":"Byte","companion_species":"capybara"}]}`, http.StatusOK)

	got, err := NewService(progress.URL, user.URL).GetLeague(context.Background(), "tok")
	require.NoError(t, err)
	require.Equal(t, []string{"u1", "u2"}, *ids)
	require.Equal(t, "Silver", got.Tier.Name)
	require.Equal(t, 2, got.MyRank)
	require.Len(t, got.Standings, 2)
	require.Equal(t, "al", got.Standings[0].Username)
	require.Equal(t, "capybara", got.Standings[0].CompanionSpecies)
	require.Equal(t, "bo", got.Standings[1].Username)
	require.True(t, got.Standings[1].IsMe)
}

func TestGetLeagueFailsClosed(t *testing.T) {
	progress, user, _ := upstreams(t, `{"error":{"code":"internal"}}`, http.StatusInternalServerError)
	_, err := NewService(progress.URL, user.URL).GetLeague(context.Background(), "tok")
	require.Error(t, err, "a user-service failure must fail the league, not show nameless rows")

	progress, user, _ = upstreams(t, `{"data":[{"id":"u1","username":"al"}]}`, http.StatusOK)
	_, err = NewService(progress.URL, user.URL).GetLeague(context.Background(), "tok")
	require.Error(t, err, "a standing with no public card must fail the league")
}

func TestGetLeagueNotJoinedSkipsUserLookup(t *testing.T) {
	progress := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{"joined":false,"tier":{"index":0,"slug":"bronze","name":"Bronze"},"standings":[]}}`))
	}))
	defer progress.Close()

	got, err := NewService(progress.URL, "http://127.0.0.1:1").GetLeague(context.Background(), "tok")
	require.NoError(t, err)
	require.False(t, got.Joined)
	require.Empty(t, got.Standings)
	require.NotNil(t, got.Standings)
}
