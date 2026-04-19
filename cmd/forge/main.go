package main

import (
	"fmt"
	"os"

	"github.com/mcchukwu/forge/internal/cli"
	"github.com/mcchukwu/forge/internal/scaffold"
)

func main() {
	opts, err := cli.ParseArgs(os.Args)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	if err := scaffold.Run(opts); err != nil {
		fmt.Printf("Error: %v", err)
	}
}
