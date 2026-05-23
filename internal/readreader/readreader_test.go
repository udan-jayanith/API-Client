package readreader_test

import (
	"Zbolt/internal/readreader"
	"io"
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestReadReader(t *testing.T) {
	// Test:
	// * write lock after calling NewReader
	// * buffering
	// * writing
	// * Content
	rr := readreader.NewReadReader(readreader.DefualtSize, nil)
	hello_world := "Hello world"
	_, err := rr.Write([]byte(hello_world))
	if err != nil {
		t.Fatalf("Writing error: %s", err.Error())
	}

	if content, err := rr.Content(); err != nil {
		t.Fatalf("Content error: %s", err.Error())
	} else if string(content) != hello_world {
		t.Fatalf("Unexpected output from call Content\nExpected '%s' but got '%s'", hello_world, string(content))
	}

	// ----------------------------------------------------------------------------
	file_path := filepath.Join("./sample", "file")
	_, err = os.Stat("./"+file_path)
	if err != nil {
		t.Fatalf("Error %s\n\nbuffer file file is missing\nRun:\n\t%s\n", err.Error(), "go run ./"+filepath.Join("./internal/readreader/sample/gen.go"))
		return
	}

	buf_file, err := os.Open(file_path)
	if err != nil {
		t.Fatalf("File error: %s", err.Error())
	}

	rr = readreader.NewReadReader(1024, nil)
	_, err = buf_file.WriteTo(rr)
	if err != nil {
		t.Fatalf("Writing error: %s", err.Error())
	}
	buf_file.Seek(0, 0)

	rr_content, err := rr.Content()
	if err != nil {
		t.Fatalf("Content error: %s", err.Error())
	}
	buf_content, err := io.ReadAll(buf_file)
	if err != nil {
		t.Fatalf("Buf file error: %s", err.Error())
	}

	if slices.Compare(rr_content, buf_content) != 0 {
		t.Fatal("Diffrent contetn found in rr_content and buf_content")
	}

	rr = readreader.NewReadReader(1024, nil)
	rr.NewReader()
	defer func() {
		if recover() == nil {
			t.Fatal("Write lock failed")
		}
	}()
	rr.Write(nil)
}

func TestReader(t *testing.T) {
	// Test:
	// * io.Reader
}
