package utf8clean

// Error types for the utf8clean processor.

import (
	"fmt"
)

// ErrInvalidCharset indicates an unsupported or invalid charset was specified.
type ErrInvalidCharset struct {
	Charset string
}

func (e ErrInvalidCharset) Error() string {
	return fmt.Sprintf("unsupported charset: %s", e.Charset)
}

func (e ErrInvalidCharset) Unwrap() error {
	return nil
}

// ErrInvalidTransferEncoding indicates an unsupported transfer encoding was specified.
type ErrInvalidTransferEncoding struct {
	Encoding string
}

func (e ErrInvalidTransferEncoding) Error() string {
	return fmt.Sprintf("unsupported transfer encoding: %s", e.Encoding)
}

func (e ErrInvalidTransferEncoding) Unwrap() error {
	return nil
}
