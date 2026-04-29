package parser

import (
	"errors"
	"regexp"

	"logprocessor/pkg/model"
)

type RegexParser struct {
	re *regexp.Regexp
}

func NewRegexParser() RegexParser {
	return RegexParser{
		re: regexp.MustCompile(`\[(.*?)\]\s+(INFO|ERROR|WARNING):\s+(.*)`),
	}
}

func (p RegexParser) Parse(line string) (model.LogEntry, error) {
	matches := p.re.FindStringSubmatch(line)

	if len(matches) != 4 {
		return model.LogEntry{}, errors.New("invalid log format")
	}

	return model.LogEntry{
		Timestamp: matches[1],
		Level:     matches[2],
		Message:   matches[3],
	}, nil
}