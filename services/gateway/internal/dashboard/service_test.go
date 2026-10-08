package dashboard

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOrderTopicsPutsFocusFirstWithoutHidingAnything(t *testing.T) {
	catalog := []TopicCard{{Slug: "sd"}, {Slug: "be"}, {Slug: "lld"}, {Slug: "dsa"}}

	got := orderTopics(catalog, []string{"dsa", "sd"})
	var slugs []string
	for _, t := range got {
		slugs = append(slugs, t.Slug)
	}
	require.Equal(t, []string{"dsa", "sd", "be", "lld"}, slugs, "focus topics in the learner's priority, then the rest in catalog order")
	require.True(t, got[0].Focused)
	require.True(t, got[1].Focused)
	require.False(t, got[2].Focused)
	require.Len(t, got, 4)

	require.Equal(t, catalog, orderTopics(catalog, nil), "no focus topics leaves the catalog order alone")
	require.Len(t, orderTopics(catalog, []string{"ghost"}), 4, "an unknown focus slug neither hides nor invents a topic")
}

func TestNextLessonIsTheCurrentNode(t *testing.T) {
	path := &pathPayload{}
	require.Nil(t, nextLesson(path))
	require.Nil(t, nextLesson(nil))

	require.NoError(t, jsonUnmarshal(`{"worlds":[{"name":"W","nodes":[
		{"label":"A","status":"done","lesson_id":"1","title":"One"},
		{"label":"B","status":"current","lesson_id":"2","title":"Two","est_minutes":4,"xp_preview":20,"in_progress":true,"kind":"lesson"},
		{"label":"C","status":"locked","lesson_id":"3","title":"Three"}]}]}`, path))

	got := nextLesson(path)
	require.NotNil(t, got)
	require.Equal(t, "2", got.LessonID)
	require.Equal(t, "Two", got.Title)
	require.Equal(t, "B", got.NodeLabel)
	require.Equal(t, "W", got.WorldName)
	require.True(t, got.InProgress)

	require.NoError(t, jsonUnmarshal(`{"worlds":[{"name":"W","nodes":[{"label":"A","status":"done"}]}]}`, path))
	require.Nil(t, nextLesson(path), "everything finished means nothing to continue")
}

func TestGetHomeAggregatesAllUpstreams(t *testing.T) {
	mux := http.NewServeMux()
	respond := func(body string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "Bearer tok", r.Header.Get("Authorization"))
			_, _ = w.Write([]byte(body))
		}
	}
	mux.HandleFunc("/api/v1/users/profile", respond(`{"data":{"onboarding_completed":true,"focus_topics":["be"],"companion":{"id":"c","name":"Kodo","species":"axolotl"}}}`))
	mux.HandleFunc("/api/v1/progress/me", respond(`{"data":{"total_xp":0,"current_level":1,"gem_balance":0,"xp_to_next_level":100}}`))
	mux.HandleFunc("/api/v1/streaks/me", respond(`{"data":{"current_streak":3,"longest_streak":5,"freeze_count":1,"streak_active_today":true}}`))
	mux.HandleFunc("/api/v1/progress/topics", respond(`{"data":[{"slug":"sd","name":"System Design","mastery":42,"skills_started":1,"skills_total":4},{"slug":"be","name":"Backend","mastery":null,"skills_started":0,"skills_total":4}]}`))
	mux.HandleFunc("/api/v1/path", func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "be", r.URL.Query().Get("focus"), "the path is ordered by the profile's focus topics")
		respond(`{"data":{"worlds":[{"name":"First Ascent","nodes":[{"label":"Why Caches Exist","status":"current","lesson_id":"L1","title":"Why Caches Exist","est_minutes":4,"xp_preview":20}]}]}}`)(w, r)
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	svc := NewService(srv.URL, srv.URL, srv.URL, srv.URL)
	home, err := svc.GetHome(context.Background(), "tok")
	require.NoError(t, err)

	require.False(t, home.OnboardingNeeded)
	require.Equal(t, 3, home.Streak.CurrentStreak)
	require.Equal(t, "Kodo", home.Companion.Name)
	require.Equal(t, []string{"be"}, home.FocusTopics)
	require.Len(t, home.Topics, 2)
	require.Equal(t, "be", home.Topics[0].Slug, "focus topic first")
	require.True(t, home.Topics[0].Focused)
	require.Nil(t, home.Topics[0].Mastery, "not started stays null, never zero")
	require.Equal(t, 42, *home.Topics[1].Mastery)
	require.NotNil(t, home.NextLesson)
	require.Equal(t, "L1", home.NextLesson.LessonID)
	require.NotEmpty(t, home.CompanionMessage)
}

func TestGetHomeFailsWhenAnyUpstreamFails(t *testing.T) {
	mux := http.NewServeMux()
	ok := func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{"data":{}}`)) }
	mux.HandleFunc("/api/v1/users/profile", ok)
	mux.HandleFunc("/api/v1/progress/me", ok)
	mux.HandleFunc("/api/v1/streaks/me", ok)
	mux.HandleFunc("/api/v1/path", ok)
	mux.HandleFunc("/api/v1/progress/topics", func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusInternalServerError) })
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	_, err := NewService(srv.URL, srv.URL, srv.URL, srv.URL).GetHome(context.Background(), "tok")
	require.Error(t, err, "the dashboard must not show partial or invented data")
}

func jsonUnmarshal(s string, v any) error { return json.Unmarshal([]byte(s), v) }
