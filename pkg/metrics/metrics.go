package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	mu           sync.Mutex
	RequestCount int
	TotalLatency time.Duration
}

var M = &Metrics{}

func RecordRequest(duration time.Duration) {
	M.mu.Lock()
	defer M.mu.Unlock()

	M.RequestCount++
	M.TotalLatency += duration
}

func GetMetrics() (int, time.Duration) {
	M.mu.Lock()
	defer M.mu.Unlock()

	return M.RequestCount, M.TotalLatency
}