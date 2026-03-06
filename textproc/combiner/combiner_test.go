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
