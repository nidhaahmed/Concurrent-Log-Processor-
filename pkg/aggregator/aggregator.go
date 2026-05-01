package aggregator

import (
	"logprocessor/pkg/model"
	"sync"
)

type Aggregator struct {
	mu     sync.Mutex
	counts map[string]int
}

func NewAggregator() *Aggregator {
	return &Aggregator{
		counts: make(map[string]int),
	}
}

func (a *Aggregator) Add(log model.LogEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.counts[log.Level]++
}

func (a *Aggregator) GetCounts() map[string]int {
	return a.counts
}