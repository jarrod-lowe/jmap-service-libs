package chunker

import (
	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor reads chunks and returns them as Chunks at fixed boundaries.
type Processor struct {
	src       textproc.ChunkProcessor
	chunkSize int
	buffer    []byte // Buffer for pulled chunks
	bufPos    int
}

// Option configures a Processor.
type Option func(*Processor)

// WithChunkSize sets the chunk size.
func WithChunkSize(n int) Option {
	return func(p *Processor) { p.chunkSize = n }
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
	return p.nextFromSource()
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
