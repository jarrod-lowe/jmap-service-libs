package textproc

// BytesProcessor reads bytes and writes processed bytes.
// Input: io.Reader, Output: []byte blocks via Next()
// Implemented by: htmlstrip, elider
type BytesProcessor interface {
	Next() ([]byte, error)
}

// ChunkProcessor reads chunks and writes chunks.
// Input: []Chunk, Output: Chunk via Next()
// Implemented by: chunker, splitter
type ChunkProcessor interface {
	Next() (Chunk, error)
}

// ChunkCombiner reads chunks and writes chunk slices.
// Input: []Chunk, Output: ChunkSlice via Next()
// Implemented by: combiner
type ChunkCombiner interface {
	Next() (ChunkSlice, error)
}
