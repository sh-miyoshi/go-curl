package option

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/pflag"
)

func getURL() (*url.URL, error) {
	args := pflag.Args()
	if len(args) == 0 {
		return nil, fmt.Errorf("no URL specified")
	} else if len(args) >= 2 {
		return nil, fmt.Errorf("Too many args. Expect 1 url but got %d", len(args))
	}

	if !strings.HasPrefix(args[0], "http://") && !strings.HasPrefix(args[0], "https://") {
		args[0] = "http://" + args[0]
	}

	return url.Parse(args[0])
}

// Init ...
func Init() (*Option, error) {
	opt := &Option{}
	help := false
	var data []string

	pflag.BoolVarP(&help, "help", "h", false, "Show help message")
	pflag.StringVarP(&opt.Method, "request", "X", "GET", "Specify request command to use")
	pflag.StringArrayVarP(&opt.Header, "header", "H", nil, "Pass custom header(s) to server")
	pflag.BoolVarP(&opt.Insecure, "insecure", "k", false, "Allow insecure server connections when using SSL")
	pflag.BoolVarP(&opt.Redirect, "location", "L", false, "Follow redirects")
	pflag.BoolVarP(&opt.Silent, "silent", "s", false, "Silent mode")
	pflag.StringVarP(&opt.Output, "output", "o", "", "Write to file instead of stdout")
	pflag.StringArrayVarP(&data, "data", "d", []string{}, "HTTP POST data")
	pflag.StringArrayVar(&opt.DataASCII, "data-ascii", []string{}, "HTTP POST ASCII data")
	pflag.StringArrayVar(&opt.DataRaw, "data-raw", []string{}, "HTTP POST data, '@' allowed")
	pflag.StringArrayVar(&opt.DataBinary, "data-binary", []string{}, "HTTP POST binary data")
	pflag.StringArrayVar(&opt.DataURL, "data-urlencode", []string{}, "HTTP POST data url encoded")

	pflag.Parse()

	if help {
		return nil, ErrHelp
	}

	for _, d := range data {
		opt.DataASCII = append(opt.DataASCII, d)
	}

	u, err := getURL()
	if err != nil {
		return nil, err
	}
	opt.URL = *u

	return opt, nil
}
