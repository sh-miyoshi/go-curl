package file

import (
	"fmt"
	"io"
	"os"
)

// RemoveReader ...
type RemoveReader struct {
	fp          *os.File
	removeChars []byte
}

// NewReader ...
func NewReader(fname string, removeChars []byte) (*RemoveReader, error) {
	fp, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	return &RemoveReader{
		fp:          fp,
		removeChars: removeChars,
	}, nil
}

// Close ...
func (r *RemoveReader) Close() error {
	if r.fp != nil {
		return r.fp.Close()
	}
	return nil
}

// Read
func (r *RemoveReader) Read(b []byte) (int, error) {
	if r.fp == nil {
		return 0, fmt.Errorf("file pointer is nil")
	}

	buf := make([]byte, len(b))
	n, err := r.fp.Read(buf)
	if err != nil && err != io.EOF {
		return n, err
	}

	// remove chars
	n = removeChar(buf[:n], r.removeChars)
	for i := 0; i < n; i++ {
		b[i] = buf[i]
	}
	return n, err
}

func removeChar(data []byte, removes []byte) int {
	wp := 0
	for rp, d := range data {
		removed := false
		for _, rc := range removes {
			if d == rc {
				removed = true
				break
			}
		}
		if !removed {
			data[wp] = data[rp]
			wp++
		}
	}

	return wp
}
