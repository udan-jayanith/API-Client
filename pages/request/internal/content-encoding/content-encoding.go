/*
 * gzip
 * deflate
 * br
 * zstd
 */
package content_encoding

import (
	"Zbolt/internal/readreader"
	"Zbolt/pages/request/internal/content-encoding/compress/br"
	"Zbolt/pages/request/internal/content-encoding/compress/deflate"
	"Zbolt/pages/request/internal/content-encoding/compress/gzip"
	"Zbolt/pages/request/internal/content-encoding/compress/zstd"
	"bufio"
	"errors"

	"io"
)

func Decode(r io.Reader, encodings []string) (*readreader.ReadReader, error) {
	closers := make([]io.ReadCloser, 0, len(encodings))
	var err error
loop:
	for _, encoding := range encodings {
		switch encoding {
		case "gzip":
			r, err = gzip.Decompress(r)
		case "deflate":
			r, err = deflate.Decompress(r)
		case "br":
			r, err = br.Decompress(r)
		case "zstd":
			r, err = zstd.Decompress(r)
		default:
			err = errors.New("Unknown encoding")
		}

		if err != nil {
			break loop
		}

		closer, ok := r.(io.ReadCloser)
		if ok {
			closers = append(closers, closer)
		}
	}

	var rr *readreader.ReadReader
	if err == nil {
		rr = readreader.NewReadReader(readreader.DefualtSize, make([]byte, 0, 2024))
		rd := bufio.NewReader(r)
		rd.WriteTo(rr)
	}

	for i := len(closers) - 1; i >= 0; i-- {
		closers[i].Close()
	}
	return rr, err
}
