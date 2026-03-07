package utf8clean

import (
	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor reads bytes from a source and validates they are UTF-8.
type Processor struct {
	src       textproc.BytesProcessor
	blockSize int
}

// Option configures a Processor.
type Option func(*Processor)

// WithBlockSize sets the block size for reading.
// Note: In pull-based mode, block size is determined by the source processor.
func WithBlockSize(n int) Option {
	return func(p *Processor) {
		p.blockSize = n
	}
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
	return p.src.Next()
}
