package http

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	nethttp "net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sh-miyoshi/go-curl/pkg/file"
	"github.com/sh-miyoshi/go-curl/pkg/option"
)

func newClient(opt *option.Option) *nethttp.Client {
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

	return client
}

func pipeRead(buf *bytes.Buffer, reader io.ReadCloser) {
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		defer reader.Close()
		if _, err := io.Copy(pw, reader); err != nil {
			fmt.Printf("Failed to read from reader: %v\n", err)
			return
		}
	}()

	buf.ReadFrom(pr)
}

func setHeaders(header nethttp.Header, rawData []string) error {
	hs := []string{}
	for _, h := range rawData {

		if h[0] == '@' {
			// read from file
			fp, err := os.Open(h[1:])
			if err != nil {
				return err
			}
			defer fp.Close()

			scanner := bufio.NewScanner(fp)
			for scanner.Scan() {
				hs = append(hs, scanner.Text())
			}
		} else {
			hs = append(hs, h)
		}
	}

	for _, h := range hs {
		d := strings.Split(h, ":")
		if len(d) == 1 {
			// check empty value
			if h[len(h)-1:] == ";" {
				header.Add(h[:len(h)-1], "")
				continue
			}
			return fmt.Errorf("Invalid header data: %s", h)
		}

		key := d[0]
		value := ""
		for i := 1; i < len(d); i++ {
			value += d[i]
		}
		header.Add(key, value)
	}

	return nil
}

func makeBody(data []option.Data) (io.Reader, error) {
	if len(data) == 0 {
		return nil, nil
	}

	buf := new(bytes.Buffer)

	for _, d := range data {
		if d.Value == "" {
			continue
		}

		writeDirect := true
		switch d.Type {
		case option.DataASCII:
			if d.Value[0] == '@' {
				writeDirect = false
				fname := d.Value[1:]
				reader, err := file.NewReader(fname, []byte{'\r', '\n'})
				if err != nil {
					return nil, err
				}
				pipeRead(buf, reader)
			}
		case option.DataRaw:
		case option.DataBinary:
			if d.Value[0] == '@' {
				writeDirect = false
				fname := d.Value[1:]
				fp, err := os.Open(fname)
				if err != nil {
					return nil, err
				}
				pipeRead(buf, fp)
			}
		case option.DataURL:
			writeDirect = false
			if index := strings.Index(d.Value, "@"); index >= 0 {
				// format @filename or name@filename
				if index != 0 {
					buf.Write([]byte(d.Value[:index]))
				}

				fname := d.Value[index+1:]
				fp, err := os.Open(fname)
				if err != nil {
					return nil, err
				}
				pipeRead(buf, fp)
			} else if index := strings.Index(d.Value, "="); index >= 0 {
				// format =content or name=content
				content := url.QueryEscape(d.Value[index+1:])
				if index != 0 {
					buf.Write([]byte(d.Value[:index+1]))
				}
				buf.Write([]byte(content))
			} else {
				// format content
				content := url.QueryEscape(d.Value)
				buf.Write([]byte(content))
			}
		}

		if writeDirect {
			buf.Write([]byte(d.Value))
		}

	}

	return buf, nil
}

func showRequest(req *nethttp.Request) {
	fmt.Printf("> %s %s %s\n", req.Method, req.URL.Path, req.Proto)
	fmt.Printf("> %s\n", req.Host)
	for key, value := range req.Header {
		fmt.Printf("> %s: %v\n", key, value)
	}
	fmt.Printf("> Content-Length: %d\n", req.ContentLength)
}

func showResponse(res *nethttp.Response) {
	fmt.Printf("< %s %d %s\n", res.Proto, res.StatusCode, res.Status)
	fmt.Printf("< Date: %s\n", time.Now().Format(time.RFC3339))
	for key, value := range res.Header {
		fmt.Printf("< %s: %v\n", key, value)
	}
	fmt.Printf("< Content-Length: %d\n", res.ContentLength)
}

// Request ...
func Request(opt *option.Option) error {
	client := newClient(opt)

	body, err := makeBody(opt.Data)
	if err != nil {
		return err
	}

	req, err := nethttp.NewRequest(opt.Method, opt.URL.String(), body)
	if err != nil {
		return err
	}
	if err := setHeaders(req.Header, opt.Header); err != nil {
		return err
	}

	if opt.Verbose {
		showRequest(req)
		fmt.Println("Request sent ...")
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if opt.Verbose {
		showResponse(res)
		fmt.Println("Got response ...")
	}

	// show result
	writer := os.Stdout
	body = res.Body

	if opt.Output != "" {
		// Output to file
		var err error
		writer, err = os.Create(opt.Output)
		if err != nil {
			return err
		}
		defer writer.Close()
		body = io.TeeReader(res.Body, newWriter(res.ContentLength, opt.Silent))
	}

	_, err = io.Copy(writer, body)
	if opt.Output != "" && !opt.Silent {
		fmt.Println("")
	}

	return err
}
