package main

import (
	"fmt"

	"github.com/sh-miyoshi/go-curl/pkg/http"
	"github.com/sh-miyoshi/go-curl/pkg/option"
	"github.com/spf13/pflag"
)

func main() {
	// TODO output

	opt, err := option.Init()
	if err != nil {
		if err == option.ErrHelp {
			pflag.PrintDefaults()
			return
		}
		fmt.Printf("Failed to parse args: %v", err)
		return
	}

	// TODO vaildate opt

	if err := http.Request(opt); err != nil {
		fmt.Printf("Failed to request server: %v", err)
	}
}
