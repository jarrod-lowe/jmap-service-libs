package utf8clean

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor reads bytes from a source and validates they are UTF-8.
type Processor struct {
	r         io.Reader               // Deprecated: use src instead
	src       textproc.BytesProcessor // NEW: pull-based source
	blockSize int
}

// Option configures a Processor.
type Option func(*Processor)

// WithBlockSize sets the block size for reading.
func WithBlockSize(n int) Option {
	return func(p *Processor) {
		p.blockSize = n
	}
}

// New creates a new Processor with the given reader and options.
// The default block size is 1024 bytes.
// Deprecated: Use NewProcessor with BytesProcessor for pull-based composition.
func New(r io.Reader, opts ...Option) *Processor {
	p := &Processor{
		r:         r,
		blockSize: 1024,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// NewProcessor creates a new Processor with the given BytesProcessor source.
// This enables pull-based lazy evaluation where the processor calls Next() on its source.
func NewProcessor(src textproc.BytesProcessor, opts ...Option) *Processor {
	p := &Processor{
		src:       src,
		blockSize: 1024,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next reads the next block of data from the source.
// Returns io.EOF when all data has been consumed.
func (p *Processor) Next() ([]byte, error) {
	// Use pull-based source if available
	if p.src != nil {
		return p.src.Next()
	}

	// Fall back to io.Reader for backward compatibility
	buf := make([]byte, p.blockSize)
	n, err := io.ReadFull(p.r, buf)

	// Handle end of file
	if err == io.EOF && n == 0 {
		return nil, io.EOF
	}

	// Handle partial read at end of file
	if err == io.ErrUnexpectedEOF {
		return buf[:n], nil
	}

	// Return the data read
	return buf[:n], err
}
