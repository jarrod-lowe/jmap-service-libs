// Package htmlstrip removes HTML markup from input while preserving text content.
//
// The processor uses a streaming HTML5 tokenizer to parse input and extract text
// content while removing all HTML tags, scripts, styles, and comments.
//
// Features:
//   - Extracts text content from HTML elements
//   - Preserves alt text from <img> tags as inline text
//   - Removes content within <script> and <style> tags
//   - Inserts newlines at block element boundaries (p, div, h1-h6, ul, ol, li, table, etc.)
//   - Handles malformed HTML gracefully using the HTML5 tokenizer
//   - Returns data in configurable blocks (default 1024 bytes)
//
// The processor is designed for memory efficiency - it does not build a DOM tree
// and processes tokens in a streaming fashion.
//
// Example:
//
//	r := strings.NewReader("<p>Hello <b>world</b></p>")
//	p := htmlstrip.NewProcessor(textprocreader.New(r))
//	result, _ := p.Next()
//	// result is []byte("Hello world")
package htmlstrip
