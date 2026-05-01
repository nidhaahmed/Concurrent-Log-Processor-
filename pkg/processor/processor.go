package processor

import (
	"logprocessor/pkg/aggregator"
	"logprocessor/pkg/parser"
	"sync"
)

func Worker(
	jobs <-chan string,
	parsers []parser.Parser,
	agg *aggregator.Aggregator,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for line := range jobs {

		var err error

		for _, p := range parsers {
			logEntry, parseErr := p.Parse(line)
			if parseErr == nil {
				agg.Add(logEntry)
				err = nil
				break
			} else {
				err = parseErr
			}
		}

		_ = err
	}
}