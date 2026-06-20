package compress

import "io"

// If error is returned io.ReadCloser must not be used.
// r must be closed.
type Decompressor = func(r io.Reader) (io.ReadCloser, error)
