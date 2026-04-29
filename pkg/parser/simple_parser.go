package parser

import (
	"errors"
	"strings"

	"logprocessor/pkg/model"
)

type SimpleParser struct{}

func (p SimpleParser) Parse(line string) (model.LogEntry, error) {
	parts := strings.SplitN(line, " ", 3)

	if len(parts) < 3 {
		return model.LogEntry{}, errors.New("invalid log format")
	}

	if parts[1] != "INFO" && parts[1] != "ERROR" && parts[1] != "WARNING" {
		return model.LogEntry{}, errors.New("invalid log level")
	}

	return model.LogEntry{
		Timestamp: parts[0],
		Level:     parts[1],
		Message:   parts[2],
	}, nil
}