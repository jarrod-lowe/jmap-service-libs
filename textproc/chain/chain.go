package chain

import (
	"io"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/chunker"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/combiner"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/elider"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/htmlstrip"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/reader"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/splitter"
	"github.com/jarrod-lowe/jmap-service-libs/textproc/utf8clean"
)

// Default configuration values
const (
	DefaultChunkSize = 4096
	DefaultMaxBytes  = 100000
	DefaultOverlap   = 2
)

// Chain composes text processors with lazy evaluation.
// The pipeline is: reader → utf8clean → htmlstrip → elider → chunker → splitter → combiner
type Chain struct {
	combiner *combiner.Processor
}

// NewReader creates a new Chain with an io.Reader as input.
// It builds the entire processor pipeline for pull-based lazy evaluation.
func NewReader(r io.Reader) (*Chain, error) {
	return NewReaderConfig(r, DefaultChunkSize, DefaultMaxBytes, DefaultOverlap)
}

// NewReaderConfig creates a new Chain with custom configuration.
func NewReaderConfig(r io.Reader, chunkSize, maxBytes, overlap int) (*Chain, error) {
	// Build the pull-based pipeline
	readerProc := reader.New(r)
	utf8Proc, err := utf8clean.NewProcessor(readerProc)
	if err != nil {
		return nil, err
	}
	htmlProc := htmlstrip.NewProcessor(utf8Proc)
	eliderProc := elider.NewProcessor(htmlProc)

	btc := &bytesToChunk{src: eliderProc}
	chunkerProc := chunker.NewProcessor(btc, chunker.WithChunkSize(chunkSize))
	splitterProc := splitter.NewProcessor(chunkerProc, maxBytes)
	combinerProc := combiner.NewProcessor(splitterProc, overlap)

	return &Chain{combiner: combinerProc}, nil
}

// Next returns the next ChunkSlice from the chain.
// It pulls data through the entire pipeline lazily.
func (c *Chain) Next() (textproc.ChunkSlice, error) {
	return c.combiner.Next()
}
