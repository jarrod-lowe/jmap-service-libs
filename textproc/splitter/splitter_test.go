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

	result, err = p.Next()
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
