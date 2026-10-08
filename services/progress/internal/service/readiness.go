package service

import (
	"github.com/prepio/prepio/services/progress/internal/store"
)

// ReadinessService computes per-topic mastery from skill scores.
type ReadinessService struct {
	readiness *store.ReadinessStore
}

// NewReadinessService creates a ReadinessService.
func NewReadinessService(readiness *store.ReadinessStore) *ReadinessService {
	return &ReadinessService{readiness: readiness}
}
