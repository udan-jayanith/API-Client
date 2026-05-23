package brotli_test

import (
	"Zbolt/internal/brotli-decoding"
	"Zbolt/internal/readreader"
	"bytes"
	"testing"

	br "github.com/molecule-man/go-brrr"
)

func TestBrotliDecoder(t *testing.T) {
	var buf bytes.Buffer
	w, err := br.NewWriter(&buf, 11)
	if err != nil {
		t.Fatal(err.Error())
	}

	content := []byte("Hello world aaaaaaaaaa ggsgs")
	_, err = w.Write(content)
	if err != nil {
		t.Fatal(err.Error())
	}
	w.Close()

	rr := readreader.NewReadReader(readreader.DefualtSize, buf.Bytes())
	rr, err = brotli.DecodeReadReader(rr, true)
	if err != nil {
		t.Fatal(err.Error())
	}

	decompressed_content, err := rr.Content()
	if err != nil {
		t.Fatal(err.Error())
	}

	if string(decompressed_content) != string(content) {
		t.Fatalf("Expected '%s' but got '%s'", content, decompressed_content)
	}
}
