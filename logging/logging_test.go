package logging

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestNew_DefaultLevel(t *testing.T) {
	// Ensure LOG_LEVEL is not set
	t.Setenv("LOG_LEVEL", "")

	var buf bytes.Buffer
	logger := New(WithOutput(&buf))

	// Debug should NOT be logged at INFO level
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("debug message should not be logged at INFO level")
	}

	// Info SHOULD be logged at INFO level
	buf.Reset()
	logger.Info("info message")
	if buf.Len() == 0 {
		t.Error("info message should be logged at INFO level")
	}
}

func TestNew_OutputsJSON(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")

	var buf bytes.Buffer
	logger := New(WithOutput(&buf))

	logger.Info("test message")

	var logEntry map[string]any
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("output is not valid JSON: %v\nOutput: %s", err, buf.String())
	}

	// Verify required fields exist
	requiredFields := []string{"time", "level", "msg"}
	for _, field := range requiredFields {
		if _, ok := logEntry[field]; !ok {
			t.Errorf("missing required field %q in JSON output", field)
		}
	}

	// Verify msg content
	if msg, ok := logEntry["msg"].(string); !ok || msg != "test message" {
		t.Errorf("expected msg='test message', got %v", logEntry["msg"])
	}

	// Verify level content
	if level, ok := logEntry["level"].(string); !ok || level != "INFO" {
		t.Errorf("expected level='INFO', got %v", logEntry["level"])
	}
}

func TestNew_EnvLevelDebug(t *testing.T) {
	t.Setenv("LOG_LEVEL", "DEBUG")

	var buf bytes.Buffer
	logger := New(WithOutput(&buf))

	// Debug SHOULD be logged when LOG_LEVEL=DEBUG
	logger.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("debug message should be logged when LOG_LEVEL=DEBUG")
	}
}

func TestNew_EnvLevelWarn(t *testing.T) {
	t.Setenv("LOG_LEVEL", "WARN")

	var buf bytes.Buffer
	logger := New(WithOutput(&buf))

	// Info should NOT be logged when LOG_LEVEL=WARN
	logger.Info("info message")
	if buf.Len() > 0 {
		t.Error("info message should not be logged when LOG_LEVEL=WARN")
	}

	// Warn SHOULD be logged
	buf.Reset()
	logger.Warn("warn message")
	if buf.Len() == 0 {
		t.Error("warn message should be logged when LOG_LEVEL=WARN")
	}
}

func TestNew_EnvLevelError(t *testing.T) {
	t.Setenv("LOG_LEVEL", "ERROR")

	var buf bytes.Buffer
	logger := New(WithOutput(&buf))

	// Warn should NOT be logged when LOG_LEVEL=ERROR
	logger.Warn("warn message")
	if buf.Len() > 0 {
		t.Error("warn message should not be logged when LOG_LEVEL=ERROR")
	}

	// Error SHOULD be logged
	buf.Reset()
	logger.Error("error message")
	if buf.Len() == 0 {
		t.Error("error message should be logged when LOG_LEVEL=ERROR")
	}
}

func TestNew_EnvLevelInvalid(t *testing.T) {
	t.Setenv("LOG_LEVEL", "INVALID")

	var buf bytes.Buffer
	logger := New(WithOutput(&buf))

	// Debug should NOT be logged (defaults to INFO)
	logger.Debug("debug message")
	if buf.Len() > 0 {
		t.Error("debug message should not be logged when LOG_LEVEL is invalid (defaults to INFO)")
	}

	// Info SHOULD be logged (defaults to INFO)
	buf.Reset()
	logger.Info("info message")
	if buf.Len() == 0 {
		t.Error("info message should be logged when LOG_LEVEL is invalid (defaults to INFO)")
	}
}

func TestNew_WithLevelOverridesEnv(t *testing.T) {
	// Set env to ERROR but override with DEBUG via option
	t.Setenv("LOG_LEVEL", "ERROR")

	var buf bytes.Buffer
	logger := New(WithLevel(slog.LevelDebug), WithOutput(&buf))

	// Debug SHOULD be logged because WithLevel overrides LOG_LEVEL
	logger.Debug("debug message")
	if buf.Len() == 0 {
		t.Error("debug message should be logged when WithLevel(DEBUG) overrides LOG_LEVEL=ERROR")
	}
}

func TestNew_WithOutput(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")

	var buf1, buf2 bytes.Buffer
	logger1 := New(WithOutput(&buf1))
	logger2 := New(WithOutput(&buf2))

	// Log to first logger
	logger1.Info("message to buf1")
	if buf1.Len() == 0 {
		t.Error("buf1 should have received output")
	}
	if buf2.Len() > 0 {
		t.Error("buf2 should not have received output from logger1")
	}

	// Log to second logger
	buf1.Reset()
	logger2.Info("message to buf2")
	if buf1.Len() > 0 {
		t.Error("buf1 should not have received output from logger2")
	}
	if buf2.Len() == 0 {
		t.Error("buf2 should have received output")
	}
}

func TestNew_DefaultOutput(t *testing.T) {
	t.Setenv("LOG_LEVEL", "")

	// Create logger without WithOutput - should not panic
	logger := New()
	if logger == nil {
		t.Error("New() should return a non-nil logger")
	}

	// Logging should work (outputs to os.Stdout, we can't easily capture this
	// without redirecting stdout, but we verify it doesn't panic)
	logger.Info("test message to stdout")
}
