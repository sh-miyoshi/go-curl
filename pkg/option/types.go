package option

import "errors"

// Option ...
type Option struct {
	Method   string
	Data     []string
	Header   []string
	Insecure bool
	Redirect bool
	Silent   bool
}

var (
	// ErrHelp ...
	ErrHelp = errors.New("show help")
)
