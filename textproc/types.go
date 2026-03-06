package textproc

// Chunk is a unit of text for embedding.
type Chunk []byte

// ChunkSlice is a slice of chunks, used by combiner output.
type ChunkSlice []Chunk
