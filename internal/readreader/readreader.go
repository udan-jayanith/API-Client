package decompressor

// See: https://pkg.go.dev/encoding/csv

import (
	"io"
	"os"
	// "compress/gzip" // Has a read closer
	// "compress/flate" // Also has a read closer
	// "github.com/andybalholm/brotli" // br // has a reader, i think this may be closed but dosen't implement io.Closer
	// "github.com/klauspost/compress/zstd" // has a reader, i think this may be closed but dosen't implement io.Closer
)

// This implements Io.ReadCloser
type ReadReader struct {
	r io.ReadCloser
	buf  []byte
	file os.File
}

// NewReader returns a new io.Reader that reads from r's internal buffer. This does not clear the internal buffer from r
func (r *ReadReader) NewReader() io.ReadCloser {
	return nil
}

// NewReadReader returns a new *ReadReader, that reads form r, if there are any compressions r get decompressed in order.
// content reads form r gets stored in a internal buffer if internal buffer size exceeds 2mb buffer will be a file in OS temporory directory.
func NewReadReader(r io.Reader, compressions []string) *ReadReader {
	return nil
}
