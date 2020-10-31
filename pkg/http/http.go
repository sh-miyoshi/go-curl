package http

import (
	"crypto/tls"
	"fmt"
	"io"
	nethttp "net/http"
	"os"
	"strings"

	"github.com/sh-miyoshi/go-curl/pkg/option"
)

// Request ...
func Request(opt *option.Option) error {
	tlsConf := tls.Config{}

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

	// TODO remove \r\n
	// TODO set body
	// get content-length
	// use io.Pipe()

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

	var body io.Reader

	// show result
	writer := os.Stdout
	body = res.Body

	if opt.Output != "" {
		// TODO /dev/null, NUL
		// Output to file
		var err error
		writer, err = os.Create(opt.Output)
		if err != nil {
			return err
		}
		defer writer.Close()
		body = io.TeeReader(res.Body, newWriter(res.ContentLength))
	}

	written, err := io.Copy(writer, body)
	fmt.Printf("\nFinished wrote %d bytes\n", written)

	return err
}
