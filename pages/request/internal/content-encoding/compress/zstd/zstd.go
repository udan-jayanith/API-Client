package zstd

import (
	"io"

	"github.com/klauspost/compress/zstd"
)

func Decompress(r io.Reader) (io.ReadCloser, error) {
	decoder, err := zstd.NewReader(r)
	if err != nil {
		return nil, err
	}
	return decoder.IOReadCloser(), nil
}
