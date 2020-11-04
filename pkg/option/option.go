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
	version := false
	var data []string
	var dataASCII []string
	var dataRaw []string
	var dataBinary []string
	var dataURL []string

	pflag.BoolVarP(&help, "help", "h", false, "Show help message")
	pflag.BoolVarP(&version, "version", "V", false, "Show version number and quit")
	pflag.StringVarP(&opt.Method, "request", "X", "", "Specify request command to use")
	pflag.StringArrayVarP(&opt.Header, "header", "H", nil, "Pass custom header(s) to server")
	pflag.BoolVarP(&opt.Insecure, "insecure", "k", false, "Allow insecure server connections when using SSL")
	pflag.BoolVarP(&opt.Redirect, "location", "L", false, "Follow redirects")
	pflag.BoolVarP(&opt.Silent, "silent", "s", false, "Silent mode")
	pflag.StringVarP(&opt.Output, "output", "o", "", "Write to file instead of stdout")
	pflag.StringArrayVarP(&data, "data", "d", []string{}, "HTTP POST data")
	pflag.StringArrayVar(&dataASCII, "data-ascii", []string{}, "HTTP POST ASCII data")
	pflag.StringArrayVar(&dataRaw, "data-raw", []string{}, "HTTP POST data, '@' allowed")
	pflag.StringArrayVar(&dataBinary, "data-binary", []string{}, "HTTP POST binary data")
	pflag.StringArrayVar(&dataURL, "data-urlencode", []string{}, "HTTP POST data url encoded")

	pflag.Parse()

	if help {
		return nil, ErrHelp
	}

	if version {
		return nil, ErrVersion
	}

	for _, d := range data {
		opt.Data = append(opt.Data, Data{Type: DataASCII, Value: d})
	}
	for _, d := range dataASCII {
		opt.Data = append(opt.Data, Data{Type: DataASCII, Value: d})
	}
	for _, d := range dataRaw {
		opt.Data = append(opt.Data, Data{Type: DataRaw, Value: d})
	}
	for _, d := range dataBinary {
		opt.Data = append(opt.Data, Data{Type: DataBinary, Value: d})
	}
	for _, d := range dataURL {
		opt.Data = append(opt.Data, Data{Type: DataURL, Value: d})
	}

	if opt.Method == "" {
		// if data is not empty, default request type is POST
		if len(opt.Data) > 0 {
			opt.Method = "POST"
		} else {
			opt.Method = "GET"
		}
	}

	u, err := getURL()
	if err != nil {
		return nil, err
	}
	opt.URL = *u

	return opt, nil
}
