package readreader

// See: https://pkg.go.dev/encoding/csv

import (
	"io"
	"os"
)

var size int = 5 * 1024 * 1024 // 5MB

// ReadReader should be closed regurdless of error.
type ReadReader struct {
	buf  []byte
	file *os.File
}

func temp_file() (*os.File, error) {
	return os.CreateTemp(os.TempDir(), "api-client-*")
}

// Write writes to the underlying bytes buffer if it's size exceed 5MB underlying buffer becomes the file.
func (r *ReadReader) Write(b []byte) (int, error) {
	if len(r.buf)+len(b) >= size {
		file, err := temp_file()
		if err != nil {
			return 0, err
		}
		r.file = file
		r.file.Write(r.buf)
		r.buf = nil
	}
	if r.file != nil {
		return r.file.Write(b)
	} else {
		r.buf = append(r.buf, b...)
		return len(b), nil
	}
}

// Close should be used to free up space
func (r *ReadReader) Close() error {
	if r.file != nil {
		// TODO: consider deleting the file.
		return r.file.Close()
	}
	r.buf = nil
	return nil
}

func (r *ReadReader) Content() ([]byte, error) {
	if r.file != nil {
		r.file.Seek(0, 0)
		return  io.ReadAll(r.file)
	}
	return r.buf, nil
}

// NewReader returns a new io.Reader that reads from r's internal buffer. This does not clear the internal buffer from r
func (r *ReadReader) NewReader() io.ReadCloser {
	panic("Not implemented")
}

// NewReadReader returns a reader that stores p, If size of p is greate then 19. p's content is stored to a file in the OS's tempory folder.
func NewReadReader(p []byte) *ReadReader {
	panic("Not implemented")
}
