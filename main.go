package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/mcchukwu/forge/internal/cli"
	"github.com/mcchukwu/forge/internal/scaffold"
)

var version = "dev"

func main() {
	opts, err := cli.ParseArgs(os.Args[1:])
	if err != nil {
		switch {
		case errors.Is(err, cli.ErrHelp):
			fmt.Print(cli.Usage)
			return
		case errors.Is(err, cli.ErrVersion):
			fmt.Printf("forge %s\n", version)
			return
		default:
			fmt.Fprintf(os.Stderr, "Error: %v\n\n%s", err, cli.Usage)
			os.Exit(1)
		}
	}

	if err := scaffold.Run(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
