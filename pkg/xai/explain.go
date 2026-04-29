package xai

import "logprocessor/pkg/model"

func Explain(log model.LogEntry) string {
	switch log.Level {
	case "ERROR":
		return "Critical issue. Immediate attention required."
	case "WARNING":
		return "Potential issue detected."
	case "INFO":
		return "System operating normally."
	default:
		return "Unknown log level."
	}
}