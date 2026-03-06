// Package textproc provides text processing and chunking functionality for email content.
//
// The processing chain flows: htmlStripper -> elider -> chunker -> splitter -> combiner.
// Each processor is memory-efficient (designed for Lambda constraints) and independently usable.
//
// Core types:
//   - Chunk: a unit of text ([]byte)
//   - ChunkSlice: a slice of chunks
//
// Processor interfaces:
//   - BytesProcessor: reads bytes via io.Reader, writes []byte blocks
//   - ChunkProcessor: reads chunks, writes chunks
//   - ChunkCombiner: reads chunks, writes chunk slices
package textproc
