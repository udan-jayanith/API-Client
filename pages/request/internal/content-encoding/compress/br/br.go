package br

import (
	"io"

	"github.com/molecule-man/go-brrr"
)

func Decompress(r io.Reader) (io.ReadCloser, error) {
	br := brrr.NewReader(r)
	return br, nil
}
