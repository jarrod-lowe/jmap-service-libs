package awsinit

import "fmt"

// TracingInitError is returned when tracing initialization fails.
type TracingInitError struct {
	Err error
}

func (e *TracingInitError) Error() string {
	return fmt.Sprintf("tracing init failed: %v", e.Err)
}

func (e *TracingInitError) Unwrap() error {
	return e.Err
}

// AWSConfigError is returned when AWS config loading fails.
type AWSConfigError struct {
	Err error
}

func (e *AWSConfigError) Error() string {
	return fmt.Sprintf("aws config failed: %v", e.Err)
}

func (e *AWSConfigError) Unwrap() error {
	return e.Err
}

// ConfigError is returned when configuration validation fails.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: %s: %s", e.Field, e.Message)
}
