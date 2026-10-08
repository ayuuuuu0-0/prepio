package service

import (
	"testing"

	"github.com/prepio/prepio/services/question/internal/store"
)

func worldsOf(nodes []store.PathNode) []string {
	var out []string
	for _, n := range nodes {
		if len(out) == 0 || out[len(out)-1] != n.WorldSlug {
			out = append(out, n.WorldSlug)
		}
	}
	return out
}

func TestOrderByFocus(t *testing.T) {
	nodes := []store.PathNode{
		{WorldSlug: "sd", WorldTopic: "system-design", NodeSlug: "sd-1"},
		{WorldSlug: "sd", WorldTopic: "system-design", NodeSlug: "sd-2"},
		{WorldSlug: "be", WorldTopic: "backend-production", NodeSlug: "be-1"},
		{WorldSlug: "misc", WorldTopic: "", NodeSlug: "misc-1"},
		{WorldSlug: "dsa", WorldTopic: "dsa-refresher", NodeSlug: "dsa-1"},
		{WorldSlug: "dsa", WorldTopic: "dsa-refresher", NodeSlug: "dsa-2"},
	}
	cases := []struct {
		name  string
		focus []string
		want  []string
	}{
		{"no focus keeps authored order", nil, []string{"sd", "be", "misc", "dsa"}},
		{"one focus topic moves its world first", []string{"dsa-refresher"}, []string{"dsa", "sd", "be", "misc"}},
		{"focus priority is kept", []string{"backend-production", "dsa-refresher"}, []string{"be", "dsa", "sd", "misc"}},
		{"an unknown topic changes nothing", []string{"cooking"}, []string{"sd", "be", "misc", "dsa"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := orderByFocus(nodes, c.focus)
			if w := worldsOf(got); !equal(w, c.want) {
				t.Fatalf("worlds = %v, want %v", w, c.want)
			}
			if len(got) != len(nodes) {
				t.Fatalf("focus must never drop nodes")
			}
		})
	}
	// Node order inside a world never changes.
	got := orderByFocus(nodes, []string{"dsa-refresher"})
	if got[0].NodeSlug != "dsa-1" || got[1].NodeSlug != "dsa-2" {
		t.Fatalf("nodes inside a world must keep their order, got %s, %s", got[0].NodeSlug, got[1].NodeSlug)
	}
	if nodes[0].NodeSlug != "sd-1" {
		t.Fatalf("orderByFocus must not modify its input")
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
