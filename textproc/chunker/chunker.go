package chunker

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor reads chunks and returns them as Chunks at fixed boundaries.
type Processor struct {
	data        [][]byte                // Deprecated: use src instead
	src         textproc.ChunkProcessor // NEW: pull-based source
	blockIndex  int
	blockOffset int
	chunkSize   int
	buffer      []byte // Buffer for pulled chunks
	bufPos      int
}

// Option configures a Processor.
type Option func(*Processor)

// WithChunkSize sets the chunk size.
func WithChunkSize(n int) Option {
	return func(p *Processor) { p.chunkSize = n }
}

// New creates a new Processor.
// Deprecated: Use NewProcessor with ChunkProcessor for pull-based composition.
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

// NewProcessor creates a new Processor with the given ChunkProcessor source.
// This enables pull-based lazy evaluation.
func NewProcessor(src textproc.ChunkProcessor, opts ...Option) *Processor {
	p := &Processor{
		src:       src,
		chunkSize: 4096,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next returns the next Chunk.
func (p *Processor) Next() (textproc.Chunk, error) {
	// Use pull-based source if available
	if p.src != nil {
		return p.nextFromSource()
	}

	// Fall back to [][]byte for backward compatibility
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

// nextFromSource pulls chunks from the source and returns fixed-size chunks.
func (p *Processor) nextFromSource() (textproc.Chunk, error) {
	// If we have buffered data, return from there
	if p.buffer != nil && p.bufPos < len(p.buffer) {
		remaining := len(p.buffer) - p.bufPos
		if remaining <= p.chunkSize {
			chunk := make(textproc.Chunk, remaining)
			copy(chunk, p.buffer[p.bufPos:])
			p.buffer = nil
			p.bufPos = 0
			return chunk, nil
		}

		chunk := make(textproc.Chunk, p.chunkSize)
		copy(chunk, p.buffer[p.bufPos:p.bufPos+p.chunkSize])
		p.bufPos += p.chunkSize
		if p.bufPos >= len(p.buffer) {
			p.buffer = nil
			p.bufPos = 0
		}
		return chunk, nil
	}

	// Need to pull more data from source
	chunk, err := p.src.Next()
	if err != nil {
		return nil, err
	}

	// If chunk fits in chunkSize, return it directly
	if len(chunk) <= p.chunkSize {
		return chunk, nil
	}

	// Chunk is too large, buffer it and return first part
	p.buffer = chunk
	p.bufPos = p.chunkSize
	result := make(textproc.Chunk, p.chunkSize)
	copy(result, chunk[:p.chunkSize])
	return result, nil
}
