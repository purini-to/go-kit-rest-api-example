package services

import (
	"context"
)

type health struct {
}

func newHealth() *health {
	return &health{}
}

func (s *health) HealthCheck(_ context.Context) string {
	return "healthy"
}
