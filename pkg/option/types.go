package option

import (
	"errors"
	"net/url"
)

// DataType ...
type DataType int

const (
	// DataASCII ...
	DataASCII DataType = iota
	// DataRaw ...
	DataRaw
	// DataBinary ...
	DataBinary
	// DataURL ...
	DataURL
)

// Data ...
type Data struct {
	Type  DataType
	Value string
}

// Option ...
type Option struct {
	Method   string
	Header   []string
	Insecure bool
	Redirect bool
	Silent   bool
	Output   string
	Data     []Data

	URL url.URL
}

var (
	// ErrHelp ...
	ErrHelp = errors.New("show help")
	// ErrVersion ...
	ErrVersion = errors.New("show version")
)
