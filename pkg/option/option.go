package option

import (
	"github.com/spf13/pflag"
)

// Init ...
func Init() (*Option, error) {
	opt := &Option{}
	help := false

	pflag.StringVarP(&opt.Method, "request", "X", "GET", "Specify request command to use")
	pflag.BoolVarP(&help, "help", "h", false, "Show help message")
	pflag.Parse()

	if help {
		return nil, ErrHelp
	}

	return opt, nil
}
