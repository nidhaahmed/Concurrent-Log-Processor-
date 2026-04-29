package parser

import "logprocessor/pkg/model"

type Parser interface {
	Parse(line string) (model.LogEntry, error)
}