package option

import "errors"

// Option ...
type Option struct {
	Method string
}

var (
	// ErrHelp ...
	ErrHelp = errors.New("show help")
)
