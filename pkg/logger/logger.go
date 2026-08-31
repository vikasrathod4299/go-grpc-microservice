package logger

/*
================================================================================
PACKAGE: pkg/logger
================================================================================

PURPOSE:
Structured logging wrapper for all microservices using Go's built-in `log/slog` package (Go 1.21+).
Structured JSON logs make debugging across microservices easy.

LEARNING GO CONCEPTS:
- Standard library `log/slog` package.
- Structured contextual logging (key-value pairs).

WHAT YOU NEED TO IMPLEMENT HERE:
1. Initialize a global or injectable `*slog.Logger`.
2. Provide helper functions: `Info(msg string, args ...any)`, `Error(...)`, `Debug(...)`.
================================================================================
*/

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func InitLogger() {
	// Initialize JSON logger outputting to stdout
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	Log = slog.New(handler)
	slog.SetDefault(Log)
}
