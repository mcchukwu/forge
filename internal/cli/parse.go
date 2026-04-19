package cli

import (
	"fmt"
)

type Options struct {
	Name    string
	Module  string
	HasMake bool
}

func ParseArgs(args []string) (Options, error) {
	if len(args) < 2 {
		return Options{}, fmt.Errorf("missing command")
	}

	if args[1] != "new" {
		return Options{}, fmt.Errorf("unknown command: %s", args[1])
	}

	if len(args) < 3 {
		return Options{}, fmt.Errorf("missing project name")
	}

	opts := Options{
		Name: args[2],
	}

	for i := 3; i < len(args); i++ {
		switch args[i] {
		case "-m", "--module":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--module missing module name")
			}

			opts.Module = args[i+1]
			i++
		case "-n", "--name":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("--name missing project name")
			}

			opts.Name = args[i+1]
			i++
		case "-M", "--make":
			opts.HasMake = true
		default:
			return opts, fmt.Errorf("unknown option: %s", args[i])
		}
	}

	return opts, nil
}
