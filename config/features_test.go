package config_test

import (
	"os"
	"testing"

	"github.com/prepio/prepio/config"
	"github.com/prepio/prepio/constants"
	"github.com/stretchr/testify/require"
)

func TestJourneyPoolSelectionDefaultOff(t *testing.T) {
	require.NoError(t, os.Unsetenv(constants.EnvJourneyPoolSelection))
	require.False(t, config.JourneyPoolSelectionEnabled())
}
