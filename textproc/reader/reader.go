package reader

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor adapts an io.Reader to the BytesProcessor interface.
// It reads data from the reader in blocks and returns them via Next().
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
// The default block size is 4096 bytes.
func New(r io.Reader, opts ...Option) *Processor {
	p := &Processor{
		r:         r,
		blockSize: 4096,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next reads the next block of data from the reader.
// Returns io.EOF when all data has been consumed.
// Implements textproc.BytesProcessor.
func (p *Processor) Next() ([]byte, error) {
	buf := make([]byte, p.blockSize)
	n, err := io.ReadFull(p.r, buf)

	// Handle end of file - no data at all
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

// Ensure Processor implements textproc.BytesProcessor
var _ textproc.BytesProcessor = (*Processor)(nil)
