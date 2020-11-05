package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sh-miyoshi/go-curl/pkg/http"
	"github.com/sh-miyoshi/go-curl/pkg/option"
	"github.com/spf13/pflag"
)

var (
	version = "latest"
)

func main() {
	// TODO output

	opt, err := option.Init()
	if err != nil {
		if err == option.ErrHelp {
			pflag.PrintDefaults()
			return
		}

		if err == option.ErrVersion {
			fmt.Printf("%s %s\n", filepath.Base(os.Args[0]), version)
			return
		}

		fmt.Printf("Failed to parse args: %v", err)
		return
	}

	if err := http.Request(opt); err != nil {
		fmt.Printf("Failed to request server: %v", err)
	}
}
