package decompressor

// See: https://pkg.go.dev/encoding/csv

import (
	"bytes"
	"io"
	"os"
	// "compress/gzip" // Has a read closer
	// "compress/flate" // Also has a read closer
	// "github.com/andybalholm/brotli" // br // has a reader
	// "github.com/klauspost/compress/zstd" // has a reader, reader should closed but dosen't implement io.Closer
)

type closer interface {
	Close()
}

// r_heandler is returned by ReadReader's NewReader. If r hasn't read fully
// this stores content from r into appropriate buffers. If buf of ReadReader is not empty NewReader must store a copy of bytes.Buffer.
// If the buufer buffer is a file pointer to the file must store as well.
type read_handler struct {
}

// This implements Io.ReadCloser
type ReadReader struct {
	r    io.Reader
	buf  *bytes.Buffer
	file os.File
}

// NewReader returns a new io.Reader that reads from r's internal buffer. This does not clear the internal buffer from r
func (r *ReadReader) NewReader() io.ReadCloser {
	return nil
}

// NewReadReader returns a new *ReadReader, that reads form r, if there are any compressions r get decompressed in order.
// content reads form r gets stored in an internal buffer. If internal buffer size exceeds 2mb buffer will be a file in OS temporory directory. If r is a io.ReadCloser ReadReader handle the closing.
func NewReadReader(r io.Reader, compressions []string) *ReadReader {
	return nil
}
