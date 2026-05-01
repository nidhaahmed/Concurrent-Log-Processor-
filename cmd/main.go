package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"logprocessor/pkg/aggregator"
	"logprocessor/pkg/parser"
	"logprocessor/pkg/processor"
	"logprocessor/pkg/reader"
)

func main() {
	lines, err := reader.ReadFile("logs.txt")
	if err != nil {
		log.Fatal(err)
	}

	// ================= DATASET AMPLIFICATION =================
	var bigLines []string
	repeat := 10000

	for i := 0; i < repeat; i++ {
		bigLines = append(bigLines, lines...)
	}
	lines = bigLines

	// ================= CONCURRENT VERSION =================

	startConcurrent := time.Now()

	jobs := make(chan string, len(lines))

	parsers := []parser.Parser{
		parser.SimpleParser{},
		parser.NewRegexParser(),
	}

	agg := aggregator.NewAggregator()

	var wg sync.WaitGroup
	numWorkers := 4

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go processor.Worker(jobs, parsers, agg, &wg)
	}

	// Send jobs
	for _, line := range lines {
		jobs <- line
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()

	durationConcurrent := time.Since(startConcurrent)

	fmt.Println("\nConcurrent Summary:")
	for level, count := range agg.GetCounts() {
		fmt.Printf("%s: %d\n", level, count)
	}

	fmt.Println("\nConcurrent Time:", durationConcurrent)

	// ================= SEQUENTIAL VERSION =================

	startSequential := time.Now()

	aggSeq := aggregator.NewAggregator()

	for _, line := range lines {
		var err error

		for _, p := range parsers {
			logEntry, parseErr := p.Parse(line)
			if parseErr == nil {
				aggSeq.Add(logEntry)
				err = nil
				break
			} else {
				err = parseErr
			}
		}

		_ = err
	}

	durationSequential := time.Since(startSequential)

	fmt.Println("\nSequential Summary:")
	for level, count := range aggSeq.GetCounts() {
		fmt.Printf("%s: %d\n", level, count)
	}

	fmt.Println("\nSequential Time:", durationSequential)
}