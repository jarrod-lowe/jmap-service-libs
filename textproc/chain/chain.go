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
	DefaultMaxBytes  = 100000
	DefaultOverlap   = 2
	DefaultByteLimit = 4000
)

// Chain composes text processors with lazy evaluation.
// The pipeline is: reader → utf8clean → htmlstrip → elider → chunker → splitter → combiner
type Chain struct {
	combiner *combiner.Processor
}

// NewReader creates a new Chain with an io.Reader as input.
// It builds the entire processor pipeline for pull-based lazy evaluation.
func NewReader(r io.Reader) (*Chain, error) {
	return NewReaderConfig(r, DefaultMaxBytes, DefaultOverlap)
}

// NewReaderConfig creates a new Chain with custom configuration.
func NewReaderConfig(r io.Reader, maxBytes, overlap int) (*Chain, error) {
	return NewReaderConfigWithByteLimit(r, maxBytes, overlap, DefaultByteLimit)
}

// NewReaderConfigWithByteLimit creates a new Chain with custom configuration including byte limit.
func NewReaderConfigWithByteLimit(r io.Reader, maxBytes, overlap, byteLimit int) (*Chain, error) {
	return NewReaderConfigWithEncoding(r, maxBytes, overlap, byteLimit, "", "")
}

// NewReaderConfigWithEncoding creates a new Chain with custom configuration including charset and transfer encoding.
// The charset and transferEncoding parameters are passed through to the utf8clean processor.
func NewReaderConfigWithEncoding(r io.Reader, maxBytes, overlap, byteLimit int, charset, transferEncoding string) (*Chain, error) {
	// Build the pull-based pipeline
	readerProc := reader.New(r)
	utf8Proc, err := utf8clean.NewProcessor(readerProc,
		utf8clean.WithCharset(charset),
		utf8clean.WithTransferEncoding(transferEncoding),
	)
	if err != nil {
		return nil, err
	}
	htmlProc := htmlstrip.NewProcessor(utf8Proc)
	eliderProc := elider.NewProcessor(htmlProc)
	chunkerProc := chunker.NewProcessor(eliderProc)
	splitterProc := splitter.NewProcessor(chunkerProc, maxBytes)
	combinerProc := combiner.NewProcessor(splitterProc, combiner.WithByteLimit(byteLimit), combiner.WithOverlap(overlap))

	return &Chain{combiner: combinerProc}, nil
}

// Next returns the next ChunkSlice from the chain.
// It pulls data through the entire pipeline lazily.
func (c *Chain) Next() (textproc.ChunkSlice, error) {
	return c.combiner.Next()
}
