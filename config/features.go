package config

import (
	"os"
	"strings"

	"github.com/prepio/prepio/constants"
)

// ReadinessV2Enabled reports whether V2 readiness is visible alongside V1.
func ReadinessV2Enabled() bool {
	return strings.EqualFold(os.Getenv(constants.EnvReadinessV2), "true")
}

// JourneyPoolSelectionEnabled reports whether journey uses pool-based question selection.
func JourneyPoolSelectionEnabled() bool {
	return strings.EqualFold(os.Getenv(constants.EnvJourneyPoolSelection), "true")
}
