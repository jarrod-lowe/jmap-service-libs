package chunker

import (
	"io"
	"regexp"
	"strings"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// Processor reads string blocks and splits them into chunks based on paragraph boundaries.
type Processor struct {
	src             textproc.StringProcessor
	buffer          string         // Input from source
	boundaryPattern *regexp.Regexp // Pre-compiled boundary detection
}

const boundaryPattern = `\n\n|\r\n\r\n|\n---\n|\n\*\*\*\n`

// NewProcessor creates a new Processor with the given StringProcessor source.
func NewProcessor(src textproc.StringProcessor) *Processor {
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
		for p.buffer == "" {
			block, err := p.src.Next()
			if err == io.EOF {
				// Source exhausted with no more data
				return "", io.EOF
			}
			if err != nil {
				// Other error
				return "", err
			}
			p.buffer += block
		}

		// Search for boundary pattern in buffer
		loc := p.boundaryPattern.FindStringIndex(p.buffer)
		if loc != nil {
			// Found boundary: extract paragraph
			paragraph := p.buffer[:loc[0]]
			p.buffer = p.buffer[loc[1]:]

			// Trim whitespace
			trimmed := strings.TrimSpace(paragraph)
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
			trimmed := strings.TrimSpace(p.buffer)
			p.buffer = ""
			if trimmed != "" {
				return textproc.Chunk(trimmed), nil
			}
			return "", io.EOF
		}
		if err != nil {
			// Other error
			return "", err
		}

		// Check if buffer would be too large
		if len(p.buffer)+len(block) > maxBufferSize {
			// Emit current buffer as-is
			trimmed := strings.TrimSpace(p.buffer)
			p.buffer = block
			if trimmed != "" {
				return textproc.Chunk(trimmed), nil
			}
			continue
		}

		// Append new block and continue looking for boundary
		p.buffer += block
	}
}
