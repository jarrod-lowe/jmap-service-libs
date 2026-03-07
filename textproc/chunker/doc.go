// Package chunker performs paragraph-based chunking.
//
// The Processor reads byte blocks from a BytesProcessor source and splits them
// into logical chunks based on paragraph boundaries. This is useful for text
// processing pipelines where content needs to be divided at natural language
// boundaries rather than arbitrary byte positions.
//
// Supported boundary types:
//   - Double newlines: \n\n
//   - Double carriage return/newlines: \r\n\r\n
//   - Horizontal rules: \n---\n and \n***\n
//
// Chunks are automatically trimmed of leading and trailing whitespace.
// Empty paragraphs (whitespace-only) are skipped.
//
// The Processor implements textproc.ChunkProcessor, making it suitable for
// pull-based lazy evaluation chains.
//
// Position in the pipeline:
//
//	After: preprocessing (htmlstrip, elider)
//	Before: splitter (which further divides chunks into embeddings)
package chunker
