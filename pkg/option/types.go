package option

import (
	"errors"
	"net/url"
)

// Option ...
type Option struct {
	Method   string
	Data     []string
	Header   []string
	Insecure bool
	Redirect bool
	Silent   bool
	Output   string

	URL url.URL
}

var (
	// ErrHelp ...
	ErrHelp = errors.New("show help")
)
