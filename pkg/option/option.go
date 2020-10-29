package option

import (
	"github.com/spf13/pflag"
)

// Init ...
func Init() (*Option, error) {
	opt := &Option{}
	help := false

	pflag.BoolVarP(&help, "help", "h", false, "Show help message")
	pflag.StringVarP(&opt.Method, "request", "X", "GET", "Specify request command to use")
	pflag.StringArrayVarP(&opt.Data, "data", "d", nil, "HTTP POST data")
	pflag.StringArrayVarP(&opt.Header, "header", "H", nil, "Pass custom header(s) to server")
	pflag.BoolVarP(&opt.Insecure, "insecure", "k", false, "Allow insecure server connections when using SSL")
	pflag.BoolVarP(&opt.Redirect, "location", "L", false, "Follow redirects")
	pflag.BoolVarP(&opt.Silent, "silent", "s", false, "Silent mode")
	pflag.Parse()

	if help {
		return nil, ErrHelp
	}

	return opt, nil
}
