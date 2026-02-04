// Package logging provides a configured slog.Logger for structured JSON logging
// in Lambda environments.
//
// It replaces the boilerplate typically found in main.go files:
//
//	// Before
//	var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
//	    Level: slog.LevelInfo,
//	}))
//
//	// After
//	var logger = logging.New()
//
// The logger reads LOG_LEVEL from the environment (DEBUG, INFO, WARN, ERROR)
// and defaults to INFO if not set or invalid.
//
// Options can override environment settings:
//
//	// Force debug logging regardless of LOG_LEVEL
//	logger := logging.New(logging.WithLevel(slog.LevelDebug))
//
//	// Capture output for testing
//	var buf bytes.Buffer
//	logger := logging.New(logging.WithOutput(&buf))
package logging
