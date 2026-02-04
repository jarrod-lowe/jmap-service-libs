package awsinit

import (
	"errors"
	"testing"
)

func TestTracingInitError_Error(t *testing.T) {
	underlying := errors.New("xray connection failed")
	err := &TracingInitError{Err: underlying}

	got := err.Error()
	want := "tracing init failed: xray connection failed"

	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestTracingInitError_Unwrap(t *testing.T) {
	underlying := errors.New("xray connection failed")
	err := &TracingInitError{Err: underlying}

	got := err.Unwrap()

	if got != underlying {
		t.Errorf("Unwrap() = %v, want %v", got, underlying)
	}
}

func TestTracingInitError_Is(t *testing.T) {
	underlying := errors.New("xray connection failed")
	err := &TracingInitError{Err: underlying}

	var target *TracingInitError
	if !errors.As(err, &target) {
		t.Error("errors.As failed for TracingInitError")
	}
}

func TestAWSConfigError_Error(t *testing.T) {
	underlying := errors.New("credentials not found")
	err := &AWSConfigError{Err: underlying}

	got := err.Error()
	want := "aws config failed: credentials not found"

	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAWSConfigError_Unwrap(t *testing.T) {
	underlying := errors.New("credentials not found")
	err := &AWSConfigError{Err: underlying}

	got := err.Unwrap()

	if got != underlying {
		t.Errorf("Unwrap() = %v, want %v", got, underlying)
	}
}

func TestAWSConfigError_Is(t *testing.T) {
	underlying := errors.New("credentials not found")
	err := &AWSConfigError{Err: underlying}

	var target *AWSConfigError
	if !errors.As(err, &target) {
		t.Error("errors.As failed for AWSConfigError")
	}
}

func TestConfigError_Error(t *testing.T) {
	err := &ConfigError{Field: "functionName", Message: "cannot be empty"}

	got := err.Error()
	want := "config error: functionName: cannot be empty"

	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestConfigError_Is(t *testing.T) {
	err := &ConfigError{Field: "functionName", Message: "cannot be empty"}

	var target *ConfigError
	if !errors.As(err, &target) {
		t.Error("errors.As failed for ConfigError")
	}
}
