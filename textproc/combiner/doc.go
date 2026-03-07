// Package combiner accumulates multiple chunks into ChunkSlices with configurable byte limits and overlap.
//
// The combiner performs a many-to-one transformation, reading chunks from a ChunkProcessor
// and returning ChunkSlices containing multiple accumulated chunks.
//
// # Configuration
//
// - ByteLimit: Maximum bytes per ChunkSlice (default: 4000)
// - Overlap: Number of chunks to overlap between successive outputs (default: 2)
//
// # Behavior
//
// Chunks are accumulated until adding the next chunk would exceed the byte limit.
// The accumulated chunks are then returned as a ChunkSlice. On the next call,
// the last N chunks (where N is the overlap setting) from the previous output
// are preserved as overlap, and additional chunks are accumulated.
//
// # Progress Guarantee
//
// If the overlap alone exceeds the byte limit, chunks are dropped from the front
// of the overlap buffer until a new chunk can be added. This ensures the combiner
// always makes progress.
//
// # Edge Cases
//
// - Single chunk exceeding byte limit: Returned as-is
// - Final chunks: All remaining chunks returned even if under limit
// - Empty source: Returns io.EOF immediately
package combiner
