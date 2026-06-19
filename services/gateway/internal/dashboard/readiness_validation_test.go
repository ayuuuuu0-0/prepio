package dashboard

import (
	"testing"

	"github.com/prepio/prepio/shared/readiness"
	"github.com/stretchr/testify/require"
)

func TestBuildComparison(t *testing.T) {
	v1 := readiness.V1Snapshot{
		Companies: []readiness.ScoreCard{
			{Company: "google", Score: 77},
			{Company: "amazon", Score: 55},
		},
		Overall: 66,
	}
	v2 := readiness.CompanyResponse{
		Companies: []readiness.CompanyEntry{
			{Company: "google", Readiness: 33},
			{Company: "amazon", Readiness: 45},
		},
		Overall: 39,
	}

	comparison := buildComparison([]string{"google", "amazon"}, v1, v2)
	require.Equal(t, 66, comparison.OverallV1)
	require.Equal(t, 39, comparison.OverallV2)
	require.Equal(t, -27, comparison.OverallDelta)
	require.Len(t, comparison.ByCompany, 2)
	require.Equal(t, -44, comparison.ByCompany[0].Delta)
}
