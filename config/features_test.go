package config_test

import (
	"os"
	"testing"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/constants"
	"github.com/stretchr/testify/require"
)

func TestReadinessV2EnabledDefaultOff(t *testing.T) {
	require.NoError(t, os.Unsetenv(constants.EnvReadinessV2))
	require.False(t, config.ReadinessV2Enabled())
}

func TestReadinessV2EnabledWhenTrue(t *testing.T) {
	require.NoError(t, os.Setenv(constants.EnvReadinessV2, "true"))
	t.Cleanup(func() { _ = os.Unsetenv(constants.EnvReadinessV2) })
	require.True(t, config.ReadinessV2Enabled())
}

func TestJourneyPoolSelectionDefaultOff(t *testing.T) {
	require.NoError(t, os.Unsetenv(constants.EnvJourneyPoolSelection))
	require.False(t, config.JourneyPoolSelectionEnabled())
}
