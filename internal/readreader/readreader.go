package readreader

// See: https://pkg.go.dev/encoding/csv

import (
	"io"
	"os"
)

// 10MB
var DefualtSize int = 10 * 1024 * 1024

// ReadReader should be closed regurdless of error.
type ReadReader struct {
	buf        []byte
	file       *os.File
	write_lock bool
	size       int
}

func temp_file() (*os.File, error) {
	return os.CreateTemp(os.TempDir(), "api-client-*")
}

// Write writes to the underlying bytes buffer if it's size exceed 5MB underlying buffer becomes the file.
// Write should not be used after calling r.NewReader()
func (r *ReadReader) Write(b []byte) (int, error) {
	if r.write_lock {
		panic("Write lock")
	}
	if len(r.buf)+len(b) >= r.size {
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
	r.write_lock = true
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
		return io.ReadAll(r.file)
	}
	return r.buf, nil
}

type reader struct {
	r  int // Read position
	rr *ReadReader
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.rr.file != nil {
		file := r.rr.file
		file.Seek(int64(r.r), 0)
		n, err = file.Read(p)
		r.r += n
		return n, err
	}

	if len(p) == 0 {
		return 0, nil
	}
	buf := r.rr.buf
	if len(buf) == r.r {
		return 0, io.EOF
	}
	src := buf[r.r:min(r.r+len(p), len(buf))]
	n = copy(p, src)
	r.r += n
	return n, nil
}

func (r *reader) Close() error {
	return nil
}

// NewReader returns a new io.Reader that reads from r's internal buffer. This does not clear the internal buffer from r.
// ReaderCloser does not closes the r.
func (r *ReadReader) NewReader() io.ReadCloser {
	r.write_lock = true
	return &reader{
		rr: r,
	}
}

/*
func (r *ReadReader) CloneReadReader() (*ReadReader, error) {
	if r.file != nil {
		rr := NewReadReader(r.size, nil)
		r := r.NewReader()
		p := make([]byte, 2048)
		for {
			n, err := r.Read(p)
			rr.Write(p[:n])
			
			if err == io.EOF {
				break
			}else if err != nil {
				return rr, err
			}
		}
	}
	// TODO: copy the r.buf
	return NewReadReader(r.size, r.buf), nil
}
 */

// NewReadReader returns a reader that stores p, If size of p is greate then 19. p's content is stored to a file in the OS's tempory folder.
func NewReadReader(size int, p []byte) *ReadReader {
	rr := ReadReader{}
	rr.size = size
	rr.Write(p)
	return &rr
}
