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

func TestWithHTMLStripper(t *testing.T) {
	r := strings.NewReader("test data")
	c := NewReader(r).WithHTMLStripper()

	if c == nil {
		t.Fatal("expected Chain to be non-nil")
	}
}

func TestWithUTF8Cleaner(t *testing.T) {
	r := strings.NewReader("test data")
	c := NewReader(r).WithUTF8Cleaner()

	if c == nil {
		t.Fatal("expected Chain to be non-nil")
	}

	if c.utf8cleaner == nil {
		t.Error("expected utf8cleaner to be non-nil")
	}
}

func TestNextReturnsChunkSlice(t *testing.T) {
	r := strings.NewReader("test data")
	c := NewReader(r).WithHTMLStripper()

	result, err := c.Next()
	if err != nil {
		t.Fatalf("expected no error on Next(), got %v", err)
	}

	if len(result) == 0 {
		t.Error("expected non-empty ChunkSlice")
	}
}

func TestNextEOF(t *testing.T) {
	r := strings.NewReader("")
	c := NewReader(r).WithHTMLStripper()

	_, err := c.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF for empty reader, got %v", err)
	}
}
