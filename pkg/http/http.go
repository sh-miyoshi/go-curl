package http

import (
	"crypto/tls"
	"fmt"
	nethttp "net/http"
	"strings"

	"github.com/sh-miyoshi/go-curl/pkg/option"
)

// Request ...
func Request(opt *option.Option) error {
	tlsConf := tls.Config{
		ServerName: opt.URL.Host,
	}

	if opt.Insecure {
		tlsConf.InsecureSkipVerify = true
	}

	tr := &nethttp.Transport{
		Proxy:           nethttp.ProxyFromEnvironment,
		TLSClientConfig: &tlsConf,
	}

	client := &nethttp.Client{
		Transport: tr,
	}

	if !opt.Redirect {
		client.CheckRedirect = func(req *nethttp.Request, via []*nethttp.Request) error {
			return nethttp.ErrUseLastResponse
		}
	}

	// TODO parse body

	req, err := nethttp.NewRequest(opt.Method, opt.URL.String(), nil)
	if err != nil {
		return err
	}
	for _, header := range opt.Header {
		d := strings.Split(header, ":")
		if len(d) != 2 {
			return fmt.Errorf("Invalid header data: %s", header)
		}
		req.Header.Add(strings.Trim(d[0], " "), strings.Trim(d[1], " "))
	}

	// TODO dump request dump, _ := httputil.DumpRequest(req, false)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// TODO show result

	return nil
}
