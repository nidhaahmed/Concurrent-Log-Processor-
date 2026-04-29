package processor

import (
	"fmt"
	
	"logprocessor/pkg/model"
	"logprocessor/pkg/parser"
)

func Worker(
	jobs <-chan string,
	results chan<- model.LogEntry,
	parsers []parser.Parser,
) {
	for line := range jobs {

		var logEntry model.LogEntry
		var err error

		fmt.Println("Worker processing:", line)

		// Try all parsers
		for _, p := range parsers {
			logEntry, err = p.Parse(line)
			if err == nil {
				results <- logEntry
				break
			}
		}

		// If no parser worked → skip silently
	}
}