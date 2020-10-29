package main

import (
	"fmt"

	"github.com/sh-miyoshi/go-curl/pkg/http"
	"github.com/sh-miyoshi/go-curl/pkg/option"
	"github.com/spf13/pflag"
)

func main() {
	opt, err := option.Init()
	if err != nil {
		if err == option.ErrHelp {
			pflag.PrintDefaults()
			return
		}
		fmt.Printf("Failed to parse args: %v", err)
		return
	}

	client := http.NewClient(opt)
	for _, u := range opt.URLs {
		client.Request(u)
	}
}
