package option

import (
	"errors"
	"net/url"
)

// Option ...
type Option struct {
	Method     string
	Header     []string
	Insecure   bool
	Redirect   bool
	Silent     bool
	Output     string
	DataASCII  []string
	DataRaw    []string
	DataBinary []string
	DataURL    []string

	URL url.URL
}

var (
	// ErrHelp ...
	ErrHelp = errors.New("show help")
)
