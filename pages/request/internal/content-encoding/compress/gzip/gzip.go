package gzip

import (
	"compress/gzip"
	"io"
)

func Decompress(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}
