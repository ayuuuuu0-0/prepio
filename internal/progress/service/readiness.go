package service

import (
	"github.com/prepio/prepio/internal/progress/store"
)

// ReadinessService computes per-topic mastery from skill scores.
type ReadinessService struct {
	readiness *store.ReadinessStore
}

// NewReadinessService creates a ReadinessService.
func NewReadinessService(readiness *store.ReadinessStore) *ReadinessService {
	return &ReadinessService{readiness: readiness}
}
