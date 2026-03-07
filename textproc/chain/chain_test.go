package chain

import (
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	r := strings.NewReader("test data")
	c := NewReader(r)

	if c == nil {
		t.Fatal("expected Chain to be non-nil")
	}
}

func TestNextReturnsChunkSlice(t *testing.T) {
	r := strings.NewReader("test data")
	c := NewReader(r)

	result, err := c.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF on Next(), got %v", err)
	}

	if len(result) == 0 {
		t.Error("expected non-empty ChunkSlice")
	}
}

func TestNextEOF(t *testing.T) {
	r := strings.NewReader("")
	c := NewReader(r)

	_, err := c.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF for empty reader, got %v", err)
	}
}

func TestFullPipelineProcessesHTML(t *testing.T) {
	// Test that the full pipeline strips HTML and returns chunks
	html := `<p>Hello <b>world</b></p><p>This is a test</p>`
	r := strings.NewReader(html)
	c := NewReader(r)

	result, err := c.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// Should have stripped HTML
	if len(result) == 0 {
		t.Error("expected non-empty ChunkSlice")
	}

	// Verify HTML was stripped (no < or > in output)
	for _, chunk := range result {
		for _, b := range chunk {
			if b == '<' || b == '>' {
				t.Errorf("HTML not stripped: found '%c' in output", b)
			}
		}
	}
}

func TestNextMultipleCalls(t *testing.T) {
	// Create a larger input to get multiple chunks
	data := strings.Repeat("test ", 1000) // ~5000 bytes
	r := strings.NewReader(data)
	c := NewReaderConfig(r, 1000, 2000, 1)

	count := 0
	for {
		_, err := c.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		count++
	}

	if count == 0 {
		t.Error("expected at least one ChunkSlice")
	}
}
