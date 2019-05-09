package services

import (
	"context"
	"github.com/purini-to/go-kit-rest-api-example/entities"
	"sync"
)

type task struct {
	mtx sync.RWMutex
	m   map[string]entities.Task
}

func newTask() *task {
	return &task{
		m: map[string]entities.Task{
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

func (s *task) GetTasks(_ context.Context) ([]entities.Task, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	tasks := make([]entities.Task, len(s.m))
	i := 0
	for _, m := range s.m {
		tasks[i] = m
		i++
	}

	return tasks, nil
}

func (s *task) GetTask(_ context.Context, id string) (entities.Task, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	p, ok := s.m[id]
	if !ok {
		return entities.Task{}, ErrNotFound
	}
	return p, nil
}
