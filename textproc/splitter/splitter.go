package splitter

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor splits chunks that exceed maxBytes into smaller chunks.
type Processor struct {
	chunks   []textproc.Chunk        // Deprecated: use src instead
	src      textproc.ChunkProcessor // NEW: pull-based source
	index    int
	maxBytes int
	buffer   textproc.Chunk
	bufPos   int
}

// Option configures a Processor.
type Option func(*Processor)

// New creates a new Processor.
// Deprecated: Use NewProcessor with ChunkProcessor for pull-based composition.
func New(chunks []textproc.Chunk, opts ...Option) *Processor {
	p := &Processor{chunks: chunks, index: 0}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// NewProcessor creates a new Processor with the given ChunkProcessor source.
// Chunks larger than maxBytes will be split into smaller chunks.
func NewProcessor(src textproc.ChunkProcessor, maxBytes int, opts ...Option) *Processor {
	p := &Processor{
		src:      src,
		maxBytes: maxBytes,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Next returns the next Chunk, splitting large chunks if necessary.
func (p *Processor) Next() (textproc.Chunk, error) {
	// Use pull-based source if available
	if p.src != nil {
		return p.nextFromSource()
	}

	// Fall back to []Chunk for backward compatibility
	if p.index >= len(p.chunks) {
		return nil, io.EOF
	}
	result := p.chunks[p.index]
	p.index++
	return result, nil
}

// nextFromSource pulls chunks from source and splits them if they exceed maxBytes.
func (p *Processor) nextFromSource() (textproc.Chunk, error) {
	// If we have buffered data, return from there
	if p.buffer != nil && p.bufPos < len(p.buffer) {
		remaining := len(p.buffer) - p.bufPos
		if remaining <= p.maxBytes {
			chunk := make(textproc.Chunk, remaining)
			copy(chunk, p.buffer[p.bufPos:])
			p.buffer = nil
			p.bufPos = 0
			return chunk, nil
		}

		chunk := make(textproc.Chunk, p.maxBytes)
		copy(chunk, p.buffer[p.bufPos:p.bufPos+p.maxBytes])
		p.bufPos += p.maxBytes
		if p.bufPos >= len(p.buffer) {
			p.buffer = nil
			p.bufPos = 0
		}
		return chunk, nil
	}

	// Pull next chunk from source
	chunk, err := p.src.Next()
	if err != nil {
		return nil, err
	}

	// If chunk fits in maxBytes, return it directly
	if len(chunk) <= p.maxBytes {
		return chunk, nil
	}

	// Chunk is too large, buffer it and return first part
	p.buffer = chunk
	p.bufPos = p.maxBytes
	result := make(textproc.Chunk, p.maxBytes)
	copy(result, chunk[:p.maxBytes])
	return result, nil
}
