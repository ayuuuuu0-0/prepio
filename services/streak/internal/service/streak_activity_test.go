package service

import (
	"testing"
	"time"

	"github.com/prepio/prepio/services/streak/internal/store"
	"github.com/stretchr/testify/require"
)

func TestApplyActivity(t *testing.T) {
	day := func(d int) time.Time { return time.Date(2026, 3, d, 0, 0, 0, 0, time.UTC) }
	dayp := func(d int) *time.Time { v := day(d); return &v }

	cases := []struct {
		name        string
		state       store.StreakState
		activity    time.Time
		freezes     int
		wantCurrent int
		wantLongest int
		wantLast    time.Time
		wantFreeze  bool
		wantChanged bool
	}{
		{"first activity starts a streak", store.StreakState{}, day(10), 0, 1, 1, day(10), false, true},
		{"same day changes nothing", store.StreakState{CurrentStreak: 3, LongestStreak: 5, LastActivityDate: dayp(10)}, day(10), 0, 3, 5, day(10), false, false},
		{"next day extends", store.StreakState{CurrentStreak: 3, LongestStreak: 3, LastActivityDate: dayp(10)}, day(11), 0, 4, 4, day(11), false, true},
		{"one missed day uses a freeze", store.StreakState{CurrentStreak: 3, LongestStreak: 7, LastActivityDate: dayp(10)}, day(12), 1, 4, 7, day(12), true, true},
		{"one missed day without freeze resets", store.StreakState{CurrentStreak: 3, LongestStreak: 7, LastActivityDate: dayp(10)}, day(12), 0, 1, 7, day(12), false, true},
		{"long gap resets", store.StreakState{CurrentStreak: 3, LongestStreak: 3, LastActivityDate: dayp(10)}, day(20), 2, 1, 3, day(20), false, true},
		{"stale earlier day changes nothing", store.StreakState{CurrentStreak: 6, LongestStreak: 6, LastActivityDate: dayp(10)}, day(9), 1, 6, 6, day(10), false, false},
		{"much older day changes nothing", store.StreakState{CurrentStreak: 6, LongestStreak: 9, LastActivityDate: dayp(10)}, day(1), 0, 6, 9, day(10), false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			state := tc.state
			got, freeze, changed := applyActivity(&state, tc.activity, tc.freezes)
			require.Equal(t, tc.wantCurrent, got.CurrentStreak)
			require.Equal(t, tc.wantLongest, got.LongestStreak)
			require.NotNil(t, got.LastActivityDate)
			require.True(t, tc.wantLast.Equal(*got.LastActivityDate))
			require.Equal(t, tc.wantFreeze, freeze)
			require.Equal(t, tc.wantChanged, changed)
		})
	}
}
