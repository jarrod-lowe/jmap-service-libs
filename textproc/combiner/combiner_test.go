package combiner

import (
	"io"
	"testing"

	"github.com/jarrod-lowe/jmap-service-libs/textproc"
)

func TestNextWrapsChunks(t *testing.T) {
	// Test that Next wraps each chunk in a ChunkSlice
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
	if len(result) != 1 {
		t.Fatalf("expected ChunkSlice length 1, got %d", len(result))
	}
	if string(result[0]) != "first" {
		t.Errorf("expected 'first', got '%s'", string(result[0]))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on second Next(), got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected ChunkSlice length 1, got %d", len(result))
	}
	if string(result[0]) != "second" {
		t.Errorf("expected 'second', got '%s'", string(result[0]))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error on third Next(), got %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected ChunkSlice length 1, got %d", len(result))
	}
	if string(result[0]) != "third" {
		t.Errorf("expected 'third', got '%s'", string(result[0]))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestNextImplementsChunkCombiner(t *testing.T) {
	// Verify Processor satisfies textproc.ChunkCombiner interface
	chunks := []textproc.Chunk{textproc.Chunk("test")}
	p := New(chunks)

	var _ textproc.ChunkCombiner = p
}

// NEW tests for pull-based composition with overlap

func TestNewProcessorCreatesProcessor(t *testing.T) {
	src := &mockChunkSource{chunks: []textproc.Chunk{[]byte("test")}}
	p := NewProcessor(src, 2)

	if p == nil {
		t.Fatal("expected Processor to be non-nil")
	}
}

func TestNewProcessorNoOverlap(t *testing.T) {
	// Test that overlapCount=0 returns just the current chunk
	src := &mockChunkSource{chunks: []textproc.Chunk{[]byte("a"), []byte("b"), []byte("c")}}
	p := NewProcessor(src, 0)

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected length 1, got %d", len(result))
	}
	if string(result[0]) != "a" {
		t.Errorf("expected 'a', got '%s'", string(result[0]))
	}

	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected length 1, got %d", len(result))
	}
	if string(result[0]) != "b" {
		t.Errorf("expected 'b', got '%s'", string(result[0]))
	}
}

func TestNewProcessorWithOverlap(t *testing.T) {
	// Test overlapCount=1 returns [previous, current]
	src := &mockChunkSource{chunks: []textproc.Chunk{[]byte("a"), []byte("b"), []byte("c")}}
	p := NewProcessor(src, 1)

	// First chunk - no previous chunks
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected length 1 for first chunk, got %d", len(result))
	}
	if string(result[0]) != "a" {
		t.Errorf("expected 'a', got '%s'", string(result[0]))
	}

	// Second chunk - should have [a, b]
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected length 2 for second chunk, got %d", len(result))
	}
	if string(result[0]) != "a" {
		t.Errorf("expected 'a', got '%s'", string(result[0]))
	}
	if string(result[1]) != "b" {
		t.Errorf("expected 'b', got '%s'", string(result[1]))
	}

	// Third chunk - should have [b, c]
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected length 2 for third chunk, got %d", len(result))
	}
	if string(result[0]) != "b" {
		t.Errorf("expected 'b', got '%s'", string(result[0]))
	}
	if string(result[1]) != "c" {
		t.Errorf("expected 'c', got '%s'", string(result[1]))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF, got %v", err)
	}
}

func TestNewProcessorWithOverlapTwo(t *testing.T) {
	// Test overlapCount=2 returns [two-previous, one-previous, current]
	src := &mockChunkSource{chunks: []textproc.Chunk{[]byte("a"), []byte("b"), []byte("c"), []byte("d")}}
	p := NewProcessor(src, 2)

	// First chunk
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 1 || string(result[0]) != "a" {
		t.Errorf("expected ['a'], got %v", result)
	}

	// Second chunk
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 || string(result[0]) != "a" || string(result[1]) != "b" {
		t.Errorf("expected ['a', 'b'], got %v", result)
	}

	// Third chunk - should have [a, b, c]
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
	if string(result[0]) != "a" || string(result[1]) != "b" || string(result[2]) != "c" {
		t.Errorf("expected ['a', 'b', 'c'], got %v", result)
	}

	// Fourth chunk - should have [b, c, d]
	result, err = p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected length 3, got %d", len(result))
	}
	if string(result[0]) != "b" || string(result[1]) != "c" || string(result[2]) != "d" {
		t.Errorf("expected ['b', 'c', 'd'], got %v", result)
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
