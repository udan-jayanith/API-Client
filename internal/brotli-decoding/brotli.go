package brotli

import (
	readreader "Zbolt/internal/readreader"
	"io"

	br "github.com/molecule-man/go-brrr"
)

func DecodeReadReader(rr *readreader.ReadReader, close_rr bool) (*readreader.ReadReader, error) {
	r := rr.NewReader()
	brr := br.NewReader(r)
	if close_rr {
		defer rr.Close()
	}
	defer r.Close()
	defer brr.Close()

	rr2 := readreader.NewReadReader(readreader.DefualtSize, nil)
	buf := make([]byte, 2048)
	for {
		n, err := brr.Read(buf)
		rr2.Write(buf[:n])
		if err != nil && err != io.EOF {
			return rr2, err
		} else if err == io.EOF {
			break
		}
	}
	return rr2, nil
}
