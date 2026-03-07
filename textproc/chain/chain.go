package chain

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/chunker"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/combiner"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/elider"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/htmlstrip"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/splitter"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/utf8clean"
)

// Chain composes text processors with lazy evaluation.
type Chain struct {
	reader       io.Reader
	utf8cleaner  *utf8clean.Processor
	htmlStripper *htmlstrip.Processor
	elider       *elider.Processor
	chunker      *chunker.Processor
	splitter     *splitter.Processor
	combiner     *combiner.Processor
	byteBlocks   [][]byte
	chunkBlocks  []textproc.Chunk
}

// NewReader creates a new Chain with an io.Reader as input.
func NewReader(r io.Reader) *Chain {
	return &Chain{reader: r}
}

// WithUTF8Cleaner adds UTF-8 cleaning to the chain.
func (c *Chain) WithUTF8Cleaner(opts ...utf8clean.Option) *Chain {
	c.utf8cleaner = utf8clean.New(c.reader, opts...)
	return c
}

// WithHTMLStripper adds HTML stripping to the chain.
func (c *Chain) WithHTMLStripper(opts ...htmlstrip.Option) *Chain {
	c.htmlStripper = htmlstrip.New(c.reader, opts...)
	return c
}

// WithElider adds content elision to the chain.
func (c *Chain) WithElider(opts ...elider.Option) *Chain {
	if c.htmlStripper != nil {
		// In a full implementation, this would use htmlStripper output
		c.elider = elider.New(c.reader, opts...)
	} else {
		c.elider = elider.New(c.reader, opts...)
	}
	return c
}

// WithChunker adds paragraph chunking to the chain.
func (c *Chain) WithChunker(opts ...chunker.Option) *Chain {
	c.chunker = chunker.New(c.byteBlocks, opts...)
	return c
}

// WithSplitter adds size-based splitting to the chain.
func (c *Chain) WithSplitter(maxBytes int, opts ...splitter.Option) *Chain {
	c.splitter = splitter.New(c.chunkBlocks, opts...)
	return c
}

// WithCombiner adds overlap combining to the chain.
func (c *Chain) WithCombiner(overlapCount int, opts ...combiner.Option) *Chain {
	c.combiner = combiner.New(c.chunkBlocks, opts...)
	return c
}

// Next returns the next ChunkSlice from the chain.
func (c *Chain) Next() (textproc.ChunkSlice, error) {
	// For the stub implementation, process through htmlStripper if present
	if c.htmlStripper != nil {
		block, err := c.htmlStripper.Next()
		if err != nil {
			return nil, err
		}
		// Wrap single block in ChunkSlice
		return textproc.ChunkSlice{textproc.Chunk(block)}, nil
	}
	return nil, io.EOF
}
