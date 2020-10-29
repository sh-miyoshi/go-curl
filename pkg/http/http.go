package http

import (
	"crypto/tls"
	nethttp "net/http"

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

	// TODO parse body

	req, err := nethttp.NewRequest(opt.Method, opt.URL.String(), nil)
	if err != nil {
		return err
	}
	// TODO add header

	// TODO dump request dump, _ := httputil.DumpRequest(req, false)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// TODO show result

	return nil
}
