package content_encoding

/*
 * gzip
 * deflate
 * br
 * zstd
 */
import (
	"Zbolt/internal/readreader"
	"Zbolt/pages/request/internal/content-encoding/compress/br"
	"Zbolt/pages/request/internal/content-encoding/compress/deflate"
	"Zbolt/pages/request/internal/content-encoding/compress/gzip"
	"Zbolt/pages/request/internal/content-encoding/compress/zstd"
	"errors"

	"io"
)

// Decode may close r
func Decode(r io.Reader, encodings []string) (readreader.ReadReader, error) {
	closers := make([]io.ReadCloser, 0, len(encodings))
	closer, ok := r.(io.ReadCloser)
	if ok {
		closers = append(closers, closer)
	}

	var err error
loop:
	for _, encoding := range encodings {
		switch encoding {
		case "gzip":
			_, err = gzip.Decompress(r)
			if err != nil {
				break loop
			}
		case "deflate":
			_, err = deflate.Decompress(r)
			if err != nil {
				break loop
			}
		case "br":
			_, err = br.Decompress(r)
			if err != nil {
				break loop
			}
		case "zstd":
			_, err = zstd.Decompress(r)
			if err != nil {
				break loop
			}
		default:
			err = errors.New("Unknown encoding")
			break loop
		}
	}

	// read from r and write it to rr.
	return readreader.ReadReader{}, err
}
