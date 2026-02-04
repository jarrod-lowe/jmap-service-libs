package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Option configures the logger.
type Option func(*config)

type config struct {
	output   io.Writer
	level    slog.Level
	levelSet bool
}

// WithLevel overrides the log level (ignores LOG_LEVEL env var).
func WithLevel(level slog.Level) Option {
	return func(c *config) {
		c.level = level
		c.levelSet = true
	}
}

// WithOutput overrides the output writer (default: os.Stdout).
func WithOutput(w io.Writer) Option {
	return func(c *config) {
		c.output = w
	}
}

// New creates a structured JSON logger for Lambda environments.
// Reads LOG_LEVEL from environment (DEBUG, INFO, WARN, ERROR).
// Defaults to INFO if not set or invalid.
func New(opts ...Option) *slog.Logger {
	cfg := &config{
		output: os.Stdout,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	level := levelFromEnv()
	if cfg.levelSet {
		level = cfg.level
	}

	return slog.New(slog.NewJSONHandler(cfg.output, &slog.HandlerOptions{
		Level: level,
	}))
}

func levelFromEnv() slog.Level {
	switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
