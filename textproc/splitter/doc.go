// Package splitter performs size-based chunking for embedding preparation.
//
// The Processor reads chunks from a ChunkProcessor source and splits them
// into smaller chunks suitable for embedding models. It tries multiple splitting
// strategies in order of preference:
//
// Supported split types:
//   - Sentence boundaries: [.!?。] followed by whitespace
//   - Word boundaries: whitespace characters (space, tab, newline)
//   - Character boundaries: UTF-8 safe byte positions
//
// Chunks are trimmed of leading and trailing whitespace.
// Empty chunks (whitespace-only) are skipped.
//
// Oversized handling:
//   - If a single content unit (no sentence/word boundaries) is within 2x the byte limit,
//     it is emitted as-is to avoid creating tiny fragments
//   - Otherwise, character splitting is forced at UTF-8 safe boundaries
//
// The Processor implements textproc.ChunkProcessor, making it suitable for
// pull-based lazy evaluation chains.
//
// Position in the pipeline:
//
//	After: chunker (paragraph-based chunking)
//	Before: combiner (recombines chunks)
package splitter
