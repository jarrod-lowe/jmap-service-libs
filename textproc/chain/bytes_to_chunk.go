package chain

import (
	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

// bytesToChunk adapts a BytesProcessor to a ChunkProcessor.
// It wraps byte slices from the source as Chunks.
type bytesToChunk struct {
	src textproc.BytesProcessor
}

// Next returns the next Chunk from the BytesProcessor.
// Implements textproc.ChunkProcessor.
func (b *bytesToChunk) Next() (textproc.Chunk, error) {
	bytes, err := b.src.Next()
	return textproc.Chunk(bytes), err
}

// Ensure bytesToChunk implements textproc.ChunkProcessor
var _ textproc.ChunkProcessor = (*bytesToChunk)(nil)
