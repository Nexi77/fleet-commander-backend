package logger

import (
	"log/slog"
	"os"
)

// Setup initializes the global slog logger based on the application environment.
func Setup(env string) {
	var handler slog.Handler

	if env == "production" || env == "staging" {
		// JSON format for AWS CloudWatch / Datadog
		opts := &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Human-readable text format for local development
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
}
