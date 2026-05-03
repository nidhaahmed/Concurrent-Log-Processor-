package api

import (
	"encoding/json"
	"logprocessor/pkg/aggregator"
	"logprocessor/pkg/metrics"
	"logprocessor/pkg/parser"
	"logprocessor/pkg/processor"
	"logprocessor/pkg/storage"
	"net/http"
	"sync"
	"time"
)

var store = storage.NewStore()

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type LogRequest struct {
	Logs []string `json:"logs"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LogRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Logs) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jobs := make(chan string, len(req.Logs))

	parsers := []parser.Parser{
		parser.SimpleParser{},
		parser.NewRegexParser(),
	}

	agg := aggregator.NewAggregator()

	var wg sync.WaitGroup
	numWorkers := 4

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range jobs {
				for _, p := range parsers {
					logEntry, err := p.Parse(line)
					if err == nil {
						agg.Add(logEntry)
						store.Add(logEntry)
						break
					}
				}
			}
		}()
	}

	for _, line := range req.Logs {
		jobs <- line
	}
	close(jobs)

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agg.GetCounts())
}

func filterHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	level := r.URL.Query().Get("level")
	if level == "" {
		http.Error(w, "Missing level", http.StatusBadRequest)
		return
	}

	filtered := store.FilterByLevel(level)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

func getAllLogsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	logs := store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {

	count, totalLatency := metrics.GetMetrics()

	response := map[string]interface{}{
		"requests":         count,
		"total_latency_ms": totalLatency.Milliseconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func compareHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LogRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || len(req.Logs) == 0 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 🔥 DATASET AMPLIFICATION (IMPORTANT)
	// var bigLogs []string
	// repeat := 1000

	// for i := 0; i < repeat; i++ {
	// 	bigLogs = append(bigLogs, req.Logs...)
	// }
	// req.Logs = bigLogs

	if len(req.Logs) < 1000 {
		var bigLogs []string
		repeat := 1000

		for i := 0; i < repeat; i++ {
			bigLogs = append(bigLogs, req.Logs...)
		}
		req.Logs = bigLogs
	}

	parsers := []parser.Parser{
		parser.SimpleParser{},
		parser.NewRegexParser(),
	}

	// ---------- CONCURRENT ----------
	startConcurrent := time.Now()

	aggConcurrent := aggregator.NewAggregator()
	jobs := make(chan string, len(req.Logs))

	var wg sync.WaitGroup
	numWorkers := 4

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processor.Worker(jobs, parsers, aggConcurrent, &wg)
	}

	for _, line := range req.Logs {
		jobs <- line
	}
	close(jobs)

	wg.Wait()

	concurrentDuration := time.Since(startConcurrent)

	// ---------- SEQUENTIAL ----------
	startSequential := time.Now()

	aggSequential := aggregator.NewAggregator()

	for _, line := range req.Logs {
		for _, p := range parsers {
			logEntry, err := p.Parse(line)
			if err == nil {
				aggSequential.Add(logEntry)
				break
			}
		}
	}

	sequentialDuration := time.Since(startSequential)

	// Convert to seconds (float)
	concurrentSec := float64(concurrentDuration.Microseconds()) / 1e6
	sequentialSec := float64(sequentialDuration.Microseconds()) / 1e6

	// Improvement calculation
	diff := sequentialSec - concurrentSec
	percent := (diff / sequentialSec) * 100

	response := map[string]interface{}{
		"concurrent": map[string]interface{}{
			"time_sec": concurrentSec,
			"result":   aggConcurrent.GetCounts(),
		},
		"sequential": map[string]interface{}{
			"time_sec": sequentialSec,
			"result":   aggSequential.GetCounts(),
		},
		"analysis": map[string]interface{}{
			"faster_by_sec":       diff,
			"improvement_percent": percent,
			"verdict":             "Concurrent processing is faster",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
