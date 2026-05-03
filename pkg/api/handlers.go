package api

import (
	"encoding/json"
	"logprocessor/pkg/aggregator"
	"logprocessor/pkg/parser"
	// "logprocessor/pkg/processor"
	"logprocessor/pkg/storage"
	"net/http"
	"sync"
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
