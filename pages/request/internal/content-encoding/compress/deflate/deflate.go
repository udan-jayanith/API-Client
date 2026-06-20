package deflate

import (
	"compress/flate"
	"io"
)

func Decompress(r io.Reader) (io.ReadCloser, error) {
	return flate.NewReader(r), nil
}
