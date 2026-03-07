package splitter

import (
	"io"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

func TestNextPassesThrough(t *testing.T) {
	// Test that Next passes chunks through unchanged
	chunks := []textproc.Chunk{
		textproc.Chunk("first"),
		textproc.Chunk("second"),
		textproc.Chunk("third"),
	}
	p := New(chunks)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "first" {
		t.Errorf("expected 'first', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if string(result) != "second" {
		t.Errorf("expected 'second', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on third Next(), got %v", err)
	}
	if string(result) != "third" {
		t.Errorf("expected 'third', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestNextImplementsChunkProcessor(t *testing.T) {
	// Verify Processor satisfies textproc.ChunkProcessor interface
	chunks := []textproc.Chunk{textproc.Chunk("test")}
	p := New(chunks)

	var _ textproc.ChunkProcessor = p
}

// NEW tests for pull-based composition with ChunkProcessor

func TestNewProcessorCreatesProcessor(t *testing.T) {
	src := &mockChunkSource{chunks: []textproc.Chunk{[]byte("test")}}
	p := NewProcessor(src, 100)

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
}

func TestNewProcessorSplitsLargeChunks(t *testing.T) {
	// Test that chunks larger than maxBytes are split
	largeChunk := make(textproc.Chunk, 100)
	for i := range largeChunk {
		largeChunk[i] = 'x'
	}
	src := &mockChunkSource{chunks: []textproc.Chunk{largeChunk}}
	p := NewProcessor(src, 30)

	// First chunk
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 30 {
		t.Errorf("expected length 30, got %d", len(result))
	}

	// Second chunk
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 30 {
		t.Errorf("expected length 30, got %d", len(result))
	}

	// Third chunk
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 30 {
		t.Errorf("expected length 30, got %d", len(result))
	}

	// Fourth chunk (remaining 10 bytes)
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 10 {
		t.Errorf("expected length 10, got %d", len(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestNewProcessorPassesThroughSmallChunks(t *testing.T) {
	// Test that small chunks pass through unchanged
	src := &mockChunkSource{chunks: []textproc.Chunk{[]byte("small"), []byte("chunks")}}
	p := NewProcessor(src, 100)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "small" {
		t.Errorf("expected 'small', got '%s'", string(result))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if string(result) != "chunks" {
		t.Errorf("expected 'chunks', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

type mockChunkSource struct {
	chunks []textproc.Chunk
	index  int
}

func (m *mockChunkSource) Next() (textproc.Chunk, error) {
	if m.index >= len(m.chunks) {
		return nil, io.EOF
	}
	result := m.chunks[m.index]
	m.index++
	return result, nil
}
