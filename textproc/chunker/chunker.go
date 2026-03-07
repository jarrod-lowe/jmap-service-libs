package chunker

import (
	"io"
	"regexp"
	"strings"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor reads byte blocks and splits them into chunks based on paragraph boundaries.
type Processor struct {
	src             textproc.BytesProcessor
	buffer          []byte         // Input from source
	boundaryPattern *regexp.Regexp // Pre-compiled boundary detection
}

const boundaryPattern = `\n\n|\r\n\r\n|\n---\n|\n\*\*\*\n`

// NewProcessor creates a new Processor with the given BytesProcessor source.
func NewProcessor(src textproc.BytesProcessor) *Processor {
	return &Processor{
		src:             src,
		boundaryPattern: regexp.MustCompile(boundaryPattern),
	}
}

// Next returns the next Chunk.
func (p *Processor) Next() (textproc.Chunk, error) {
	const maxBufferSize = 1 << 20 // 1MB

	for {
		// Read more data if buffer is empty
		for len(p.buffer) == 0 {
			block, err := p.src.Next()
			if err == io.EOF {
				// Source exhausted with no more data
				return nil, io.EOF
			}
			if err != nil {
				// Other error
				return nil, err
			}
			p.buffer = append(p.buffer, block...)
		}

		// Search for boundary pattern in buffer
		loc := p.boundaryPattern.FindIndex(p.buffer)
		if loc != nil {
			// Found boundary: extract paragraph
			paragraph := p.buffer[:loc[0]]
			p.buffer = p.buffer[loc[1]:]

			// Trim whitespace
			trimmed := strings.TrimSpace(string(paragraph))
			if trimmed != "" {
				return textproc.Chunk(trimmed), nil
			}
			// Empty paragraph, continue to next
			continue
		}

		// No boundary found - try to read more data
		block, err := p.src.Next()
		if err == io.EOF {
			// Source exhausted, emit remaining content
			trimmed := strings.TrimSpace(string(p.buffer))
			p.buffer = nil
			if trimmed != "" {
				return textproc.Chunk(trimmed), nil
			}
			return nil, io.EOF
		}
		if err != nil {
			// Other error
			return nil, err
		}

		// Check if buffer would be too large
		if len(p.buffer)+len(block) > maxBufferSize {
			// Emit current buffer as-is
			trimmed := strings.TrimSpace(string(p.buffer))
			p.buffer = block
			if trimmed != "" {
				return textproc.Chunk(trimmed), nil
			}
			continue
		}

		// Append new block and continue looking for boundary
		p.buffer = append(p.buffer, block...)
	}
}
