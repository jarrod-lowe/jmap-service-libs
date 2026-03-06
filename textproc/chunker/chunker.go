package chunker

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor reads byte blocks and returns them as Chunks at fixed boundaries.
type Processor struct {
	data        [][]byte
	blockIndex  int
	blockOffset int
	chunkSize   int
}

// Option configures a Processor.
type Option func(*Processor)

// WithChunkSize sets the chunk size.
func WithChunkSize(n int) Option {
	return func(p *Processor) { p.chunkSize = n }
}

// New creates a new Processor.
func New(data [][]byte, opts ...Option) *Processor {
	p := &Processor{
		data:        data,
		blockIndex:  0,
		blockOffset: 0,
		chunkSize:   4096,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next returns the next Chunk.
func (p *Processor) Next() (textproc.Chunk, error) {
	if p.blockIndex >= len(p.data) {
		return nil, io.EOF
	}

	currentBlock := p.data[p.blockIndex]
	remainingInBlock := len(currentBlock) - p.blockOffset

	// If remaining data fits in one chunk, return it all
	if remainingInBlock <= p.chunkSize {
		chunk := make(textproc.Chunk, remainingInBlock)
		copy(chunk, currentBlock[p.blockOffset:])
		p.blockIndex++
		p.blockOffset = 0
		return chunk, nil
	}

	// Return chunkSize bytes
	chunk := make(textproc.Chunk, p.chunkSize)
	copy(chunk, currentBlock[p.blockOffset:p.blockOffset+p.chunkSize])
	p.blockOffset += p.chunkSize

	// Move to next block if current block exhausted
	if p.blockOffset >= len(currentBlock) {
		p.blockIndex++
		p.blockOffset = 0
	}

	return chunk, nil
}
