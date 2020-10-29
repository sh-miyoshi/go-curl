package http

import (
	nethttp "net/http"
	neturl "net/url"

	"github.com/sh-miyoshi/go-curl/pkg/option"
)

// Client ...
type Client struct {
	client *nethttp.Client
}

// NewClient ...
func NewClient(opt *option.Option) *Client {
	return &Client{}
}

// Request ...
func (c *Client) Request(url neturl.URL) {

}
