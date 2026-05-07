package decompressor

// See: https://pkg.go.dev/encoding/csv

import (
	"bytes"
	"io"
	"os"
	// "compress/gzip" // Has a read closer
	// "compress/flate" // Also has a read closer
	// "github.com/andybalholm/brotli" // br // has a reader
	// "github.com/klauspost/compress/zstd" // has a reader, reader should be closed but dosen't implement io.Closer
)

type closer interface {
	Close()
}

// r_heandler is returned by ReadReader's NewReader. If r hasn't read fully
// this stores content from r into appropriate buffers. If buf of ReadReader is not empty NewReader must store a copy of bytes.Buffer.
// If the buufer buffer is a file pointer to the file must store as well.
type read_handler struct {
}

var count int = 0

// This implements Io.ReadCloser
type ReadReader struct {
	id   int
	buf  *bytes.Buffer
	file *os.File
}

func (r *ReadReader) Id() int {
	return 0
}

// Decode uncompreses the underlying buffers content using given compressions formats.
func (r *ReadReader) Decode(compressions []string) error {
	return nil
}

// NewReader returns a new io.Reader that reads from r's internal buffer. This does not clear the internal buffer from r
func (r *ReadReader) NewReader() io.ReadCloser {
	return nil
}

// Write writes to the underlying bytes buffer if it's size exceed 2m underlying buffer becomes the file.
func (r *ReadReader) Write(p []byte) {

}

// Close should be used to free up spacce
func (r *ReadReader) Close() error {
	return nil
}

// NewReadReader returns a reader that stores p, If size of p is greate then 19. p's content is stored to a file in the OS's tempory folder.
func NewReadReader(p []byte) *ReadReader {
	return nil
}
