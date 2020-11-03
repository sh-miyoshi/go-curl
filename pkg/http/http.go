package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	nethttp "net/http"
	"os"
	"strings"

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

func makeBody(dataASCII, dataRaw, dataBinary, dataURL []string) (io.Reader, error) {
	if len(dataASCII) == 0 && len(dataRaw) == 0 && len(dataBinary) == 0 && len(dataURL) == 0 {
		return nil, nil
	}

	buf := new(bytes.Buffer)
	//var mu sync.Mutex

	for _, d := range dataASCII {
		if d == "" {
			continue
		}
		if d[0] == '@' {
			fname := d[1:]
			reader, err := file.NewReader(fname, []byte{'\r', '\n'})
			if err != nil {
				return nil, err
			}
			pr, pw := io.Pipe()

			go func() {
				defer pw.Close()
				defer reader.Close()
				if _, err := io.Copy(pw, reader); err != nil {
					fmt.Printf("Failed to read %s: %v\n", fname, err)
					return
				}
			}()

			buf.ReadFrom(pr)
		} else {
			buf.Write([]byte(d))
		}
	}

	for _, d := range dataRaw {
		if d == "" {
			continue
		}
		buf.Write([]byte(d))
	}

	for _, d := range dataBinary {
		if d == "" {
			continue
		}
		if d[0] == '@' {
			fname := d[1:]
			fp, err := os.Open(fname)
			if err != nil {
				return nil, err
			}
			pr, pw := io.Pipe()

			go func() {
				defer pw.Close()
				defer fp.Close()
				if _, err := io.Copy(pw, fp); err != nil {
					fmt.Printf("Failed to read %s: %v\n", fname, err)
					return
				}
			}()

			buf.ReadFrom(pr)
		} else {
			buf.Write([]byte(d))
		}
	}

	for _, d := range dataURL {
		if d == "" {
			continue
		}
		// TODO implement this
	}

	return buf, nil
}

// Request ...
func Request(opt *option.Option) error {
	client := newClient(opt)

	body, err := makeBody(opt.DataASCII, opt.DataRaw, opt.DataBinary, opt.DataURL)
	if err != nil {
		return err
	}

	req, err := nethttp.NewRequest(opt.Method, opt.URL.String(), body)
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
		body = io.TeeReader(res.Body, newWriter(res.ContentLength, opt.Silent))
	}

	_, err = io.Copy(writer, body)
	if opt.Output != "" && !opt.Silent {
		fmt.Println("")
	}

	return err
}
