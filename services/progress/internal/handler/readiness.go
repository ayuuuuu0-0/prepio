package handler

import (
	"github.com/prepio/prepio/services/progress/internal/service"
)

// ReadinessHandler serves topic mastery endpoints.
type ReadinessHandler struct {
	readiness *service.ReadinessService
}

// NewReadinessHandler creates a ReadinessHandler.
func NewReadinessHandler(readiness *service.ReadinessService) *ReadinessHandler {
	return &ReadinessHandler{readiness: readiness}
}
