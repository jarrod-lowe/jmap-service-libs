package chain

import (
	"io"
	"strings"
	"testing"
)

func TestNewReader(t *testing.T) {
	r := strings.NewReader("test data")
	c, err := NewReader(r)
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	if c == nil {
		t.Fatal("expected Chain to be non-nil")
	}
}

func TestNextReturnsChunkSlice(t *testing.T) {
	r := strings.NewReader("test data")
	c, err := NewReader(r)
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

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
	c, err := NewReader(r)
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

	_, err = c.Next()
	if err != io.EOF {
		t.Errorf("expected io.EOF for empty reader, got %v", err)
	}
}

func TestFullPipelineProcessesHTML(t *testing.T) {
	// Test that the full pipeline strips HTML and returns chunks
	html := `<p>Hello <b>world</b></p><p>This is a test</p>`
	r := strings.NewReader(html)
	c, err := NewReader(r)
	if err != nil {
		t.Fatalf("NewReader failed: %v", err)
	}

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

func TestParagraphBasedChunking(t *testing.T) {
	// Test that paragraph-based chunking works
	data := "Paragraph one\n\nParagraph two\n\nParagraph three"
	r := strings.NewReader(data)
	c, err := NewReaderConfig(r, 1000, 1)
	if err != nil {
		t.Fatalf("NewReaderConfig failed: %v", err)
	}

	result, err := c.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// Should have at least one chunk with paragraph content
	if len(result) == 0 {
		t.Error("expected non-empty ChunkSlice")
	}

	// Verify paragraphs are present
	content := strings.Join(strings.Fields(string(result[0])), " ")
	if content == "" {
		t.Error("expected paragraph content")
	}
}

func TestNextMultipleCalls(t *testing.T) {
	// Create input with multiple paragraphs
	data := "Paragraph one\n\nParagraph two\n\nParagraph three\n\nParagraph four"
	r := strings.NewReader(data)
	c, err := NewReaderConfig(r, 1000, 1)
	if err != nil {
		t.Fatalf("NewReaderConfig failed: %v", err)
	}

	count := 0
	for {
		_, err = c.Next()
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

func TestNewReaderConfigWithEncoding(t *testing.T) {
	data := "test data"
	r := strings.NewReader(data)
	c, err := NewReaderConfigWithEncoding(r, 1000, 1, 4000, "", "")
	if err != nil {
		t.Fatalf("NewReaderConfigWithEncoding failed: %v", err)
	}

	if c == nil {
		t.Fatal("expected Chain to be non-nil")
	}

	// Verify chain can process data
	result, err := c.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF on Next(), got %v", err)
	}

	if len(result) == 0 {
		t.Error("expected non-empty ChunkSlice")
	}
}

func TestNewReaderConfigWithEncodingInvalidCharset(t *testing.T) {
	data := "test data"
	r := strings.NewReader(data)
	_, err := NewReaderConfigWithEncoding(r, 1000, 1, 4000, "invalid-charset-xyz", "")
	if err == nil {
		t.Error("expected error for invalid charset")
	}
}

func TestNewReaderConfigWithEncodingInvalidTransferEncoding(t *testing.T) {
	data := "test data"
	r := strings.NewReader(data)
	_, err := NewReaderConfigWithEncoding(r, 1000, 1, 4000, "", "invalid-encoding")
	if err == nil {
		t.Error("expected error for invalid transfer encoding")
	}
}
