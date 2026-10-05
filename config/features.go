package config

import (
	"os"
	"strings"

	"github.com/prepio/prepio/constants"
)

// JourneyPoolSelectionEnabled reports whether journey uses pool-based question selection.
func JourneyPoolSelectionEnabled() bool {
	return strings.EqualFold(os.Getenv(constants.EnvJourneyPoolSelection), "true")
}
