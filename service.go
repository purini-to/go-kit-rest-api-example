package go_kit_rest_api_example

import (
	"context"
	"errors"
	"sync"
)

type Service interface {
	GetTasks(ctx context.Context) ([]Task, error)
	GetTask(_ context.Context, id string) (Task, error)
}

type Task struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	State string `json:"state"`
}

var (
	ErrNotFound = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Task
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]Task{
			"5914b29abcbe09ab9478d156": {
				ID:    "5914b29abcbe09ab9478d156",
				Title: "Custom Fields in CSV exports",
				State: "complete",
			},
			"5914b2a1ee4421b88542b59a": {
				ID:    "5914b2a1ee4421b88542b59a",
				Title: "Custom Fields on mobile",
				State: "incomplete",
			},
		},
	}
}

func (s *inmemService) GetTasks(_ context.Context) ([]Task, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	tasks := make([]Task, len(s.m))
	i := 0
	for _, m := range s.m {
		tasks[i] = m
		i++
	}

	return tasks, nil
}

func (s *inmemService) GetTask(_ context.Context, id string) (Task, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	p, ok := s.m[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return p, nil
}
