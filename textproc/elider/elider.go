package elider

import (
	"io"
)

// Processor reads bytes from an io.Reader and returns them in blocks.
// For the initial stub implementation, it reads blocks and passes them through unmodified.
type Processor struct {
	r         io.Reader
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

// Next reads the next block of data from the reader.
// Returns io.EOF when all data has been consumed.
func (p *Processor) Next() ([]byte, error) {
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
