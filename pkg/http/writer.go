package http

import (
	"fmt"
	"io"
)

type progressWriter struct {
	current int64
	total   int64
}

// Write implements the io.Writer interface.
// Always completes and never returns an error.
func (w *progressWriter) Write(p []byte) (n int, e error) {
	n = len(p)
	w.current += int64(n)
	// TODO silent mode
	percent := float64(w.current) * 100 / float64(w.total)
	fmt.Printf("\rReceived %d bytes in %d (%d%%)", w.current, w.total, int(percent))
	return
}

func newWriter(total int64) io.Writer {
	return &progressWriter{total: total}
}
