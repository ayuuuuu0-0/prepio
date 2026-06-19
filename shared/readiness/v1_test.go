package readiness_test

import (
	"testing"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/shared/readiness"
	"github.com/stretchr/testify/require"
)

func TestComputeV1CompanyScore(t *testing.T) {
	score := readiness.ComputeV1CompanyScore(readiness.CompanyStats{
		Company:  "google",
		Answered: 4,
		Correct:  3,
		ScoreAvg: 80,
	})
	require.Equal(t, (75+80)/2, score)
}

func TestComputeV1CompanyScoreCapsAtMax(t *testing.T) {
	score := readiness.ComputeV1CompanyScore(readiness.CompanyStats{
		Company:  "google",
		Answered: 10,
		Correct:  10,
		ScoreAvg: 100,
	})
	require.Equal(t, config.MaxCompanyReadiness, score)
}

func TestComputeV1Overall(t *testing.T) {
	overall := readiness.ComputeV1Overall([]readiness.CompanyScore{
		{Company: "google", Score: 70},
		{Company: "amazon", Score: 50},
	})
	require.Equal(t, 60, overall)
}
