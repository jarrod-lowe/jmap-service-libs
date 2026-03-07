package htmlstrip

import (
	"io"
	"strings"
	"testing"
)

func TestNewProcessor(t *testing.T) {
	// Test that New creates a processor with default block size
	r := strings.NewReader("test data")
	p := New(r)

	if p == nil {
		t.Fatal("expected processor to be non-nil")
	}

	// Verify default block size is 1024
	if p.blockSize != 1024 {
		t.Errorf("expected default blockSize 1024, got %d", p.blockSize)
	}
}

func TestNewProcessorWithOptions(t *testing.T) {
	// Test that New with options sets custom block size
	r := strings.NewReader("test data")
	p := New(r, WithBlockSize(512))

	if p == nil {
		t.Fatal("expected processor to be non-nil")
	}

	if p.blockSize != 512 {
		t.Errorf("expected blockSize 512, got %d", p.blockSize)
	}
}

func TestNextSingleBlock(t *testing.T) {
	// Test reading plain text (passthrough)
	data := "hello world"
	r := strings.NewReader(data)
	p := New(r, WithBlockSize(1024))

	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(result) != data {
		t.Errorf("expected '%s', got '%s'", data, string(result))
	}
}

func TestNextEmptyReader(t *testing.T) {
	// Test reading from an empty reader
	r := strings.NewReader("")
	p := New(r)

	result, err := p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF immediately, got %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result with EOF, got %v", result)
	}
}

func TestNextEOFThenNext(t *testing.T) {
	// Test that Next continues to return EOF after first EOF
	r := strings.NewReader("test")
	p := New(r, WithBlockSize(10))

	// First call should succeed
	result, err := p.Next()
	if err != nil {
		t.Fatalf("expected no error on first Next(), got %v", err)
	}
	if string(result) != "test" {
		t.Errorf("expected 'test', got '%s'", string(result))
	}

	// Second call should return EOF
	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}

	// Third call should also return EOF
	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on third Next(), got %v", err)
	}
}

func TestStripBasicHTML(t *testing.T) {
	// Test basic HTML stripping: <p>Hello <b>world</b></p> should produce "Hello world"
	data := `<p>Hello <b>world</b></p>`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(result) != "Hello world" {
		t.Errorf("expected 'Hello world', got '%s'", string(result))
	}

	// Verify we've consumed all input
	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestImgAltText(t *testing.T) {
	// Test that img alt text is preserved: <img src="x.jpg" alt="Photo"> should produce "Photo"
	data := `<img src="x.jpg" alt="Photo">`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(result) != "Photo" {
		t.Errorf("expected 'Photo', got '%s'", string(result))
	}

	// Verify we've consumed all input
	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestScriptStyleRemoval(t *testing.T) {
	// Test that content inside script and style tags is removed
	// Note: <p> tags are block elements that insert newlines
	data := `<p>Hello</p><script>alert('bad');</script><p>World</p>`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	expected := "Hello\nWorld"
	if string(result) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(result))
	}

	// Verify we've consumed all input
	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestBlockElementSpacing(t *testing.T) {
	// Test that block elements insert newlines: <p>Para1</p><p>Para2</p> should produce "Para1\nPara2"
	data := `<p>Para1</p><p>Para2</p>`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	expected := "Para1\nPara2"
	if string(result) != expected {
		t.Errorf("expected '%s', got '%s'", expected, string(result))
	}

	// Verify we've consumed all input
	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestLinkHandling(t *testing.T) {
	// Test that link text is extracted but href is ignored: <a href="url">text</a> should produce "text"
	data := `<a href="https://example.com">Click here</a>`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(result) != "Click here" {
		t.Errorf("expected 'Click here', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestMalformedHTML(t *testing.T) {
	// Test that malformed HTML still produces output: <b>unclosed should still work
	data := `<b>unclosed text`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(result) != "unclosed text" {
		t.Errorf("expected 'unclosed text', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestMixedContent(t *testing.T) {
	// Test mixed content: "Hello <p>World</p>" should produce "Hello World"
	data := `Hello <p>World</p>`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(result) != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestAngleBracketsNotTags(t *testing.T) {
	// Test that angle brackets that aren't tags are preserved: "Price: $5 < $10"
	data := `Price: $5 < $10`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	if string(result) != "Price: $5 < $10" {
		t.Errorf("expected 'Price: $5 < $10', got '%s'", string(result))
	}

	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}

func TestPartialTags(t *testing.T) {
	// Test that partial/unclosed tags are handled gracefully without panicking
	// The HTML tokenizer treats "<b text" as an incomplete tag
	data := `some <b text`
	r := strings.NewReader(data)
	p := New(r)

	result, err := p.Next()
	if err != nil && err != io.EOF {
		t.Fatalf("expected no error or EOF, got %v", err)
	}

	// The tokenizer handles this - output may vary based on tokenizer implementation
	// Just verify we got something without error
	if len(result) == 0 {
		t.Errorf("expected some output, got empty result")
	}

	// Verify we can call Next again without error
	_, err = p.Next()
	if err != io.EOF {
		t.Fatalf("expected io.EOF on second Next(), got %v", err)
	}
}
