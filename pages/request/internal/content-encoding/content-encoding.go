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

	"io"
)

func Decode(r io.Reader, encodings []string) (readreader.ReadReader, error) {
	//rr := readreader.NewReadReader(readreader.DefualtSize, make([]byte, 0, 2024))
	for _, encoding := range encodings {
		switch encoding {
		case "gzip":
			_, err := gzip.Decompress(r)
			if err != nil {
				return readreader.ReadReader{}, err
			}
		case "deflate":
			_, err := deflate.Decompress(r)
			if err != nil {
				return readreader.ReadReader{}, err
			}
		case "br":
			_, err := br.Decompress(r)
			if err != nil {
				return readreader.ReadReader{}, err
			}
		case "zstd":
			_, err := zstd.Decompress(r)
			if err != nil {
				return readreader.ReadReader{}, err
			}
		}
	}
	return readreader.ReadReader{}, nil
}
