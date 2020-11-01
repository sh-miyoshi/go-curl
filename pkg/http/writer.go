package http

import (
	"fmt"
	"io"
)

type progressWriter struct {
	current int64
	total   int64
	silent  bool
}

// Write implements the io.Writer interface.
// Always completes and never returns an error.
func (w *progressWriter) Write(p []byte) (n int, e error) {
	n = len(p)
	w.current += int64(n)
	if !w.silent {
		percent := float64(w.current) * 100 / float64(w.total)
		fmt.Printf("\rReceived %d bytes in %d (%d%%)", w.current, w.total, int(percent))
	}
	return
}

func newWriter(total int64, silent bool) io.Writer {
	return &progressWriter{
		total:  total,
		silent: silent,
	}
}
