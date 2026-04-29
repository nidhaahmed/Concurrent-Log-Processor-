package main

import (
	"fmt"
	"log"

	"logprocessor/pkg/model"
	"logprocessor/pkg/parser"
	"logprocessor/pkg/processor"
	"logprocessor/pkg/reader"
	"logprocessor/pkg/xai"
)

func main() {
	lines, err := reader.ReadFile("logs.txt")
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan string, len(lines))
	results := make(chan model.LogEntry, len(lines))

	// Parsers
	parsers := []parser.Parser{
		parser.SimpleParser{},
		parser.NewRegexParser(),
	}

	// Start workers
	numWorkers := 3
	for i := 0; i < numWorkers; i++ {
		go processor.Worker(jobs, results, parsers)
	}

	// Send jobs
	for _, line := range lines {
		jobs <- line
	}
	close(jobs)

	// Collect results
	for i := 0; i < len(lines); i++ {
		logEntry := <-results
		fmt.Println(logEntry)
		fmt.Println("Explanation:", xai.Explain(logEntry))
	}
}