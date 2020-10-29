package main

import (
	"fmt"

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
	}

	fmt.Printf("opt: %v\n", opt)
}
