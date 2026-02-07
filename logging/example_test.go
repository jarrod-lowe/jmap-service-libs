package logging_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jarrod-lowe/jmap-service-libs/logging"
)

func ExampleNew() {
	var buf bytes.Buffer
	logger := logging.New(
		logging.WithLevel(slog.LevelInfo),
		logging.WithOutput(&buf),
	)
	logger.Info("hello world")
	// Verify JSON output contains the message
	fmt.Println(strings.Contains(buf.String(), `"msg":"hello world"`))
	// Output: true
}

func ExampleWithLevel() {
	var buf bytes.Buffer
	logger := logging.New(
		logging.WithLevel(slog.LevelError),
		logging.WithOutput(&buf),
	)
	// Info is below Error level, so this won't be logged
	logger.Info("not logged")
	fmt.Println(buf.Len() == 0)
	// Output: true
}
