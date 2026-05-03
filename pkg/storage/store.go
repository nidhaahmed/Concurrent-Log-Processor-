package storage

import (
	"logprocessor/pkg/model"
	"sync"
)

type Store struct {
	mu   sync.Mutex
	logs []model.LogEntry
}

func NewStore() *Store {
	return &Store{
		logs: []model.LogEntry{},
	}
}

func (s *Store) Add(log model.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logs = append(s.logs, log)
}

func (s *Store) GetAll() []model.LogEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make([]model.LogEntry, len(s.logs))
	copy(result, s.logs)

	return result
}

func (s *Store) FilterByLevel(level string) []model.LogEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []model.LogEntry
	for _, log := range s.logs {
		if log.Level == level {
			result = append(result, log)
		}
	}
	return result
}