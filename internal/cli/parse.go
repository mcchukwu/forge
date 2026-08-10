package cli

import (
	"errors"
	"fmt"
	"strings"
)

const Usage = `forge - scaffold new Go projects

Usage:
  forge new <name> [flags]
  forge help
  forge version

Commands:
  new       Create a new Go project
  help      Show this help message
  version   Show the forge version

Flags:
  -m, --module <name>   Go module name (default: the project name)
  -n, --name <name>     Override the project name
  -M, --make            Also generate a Makefile
  -f, --force           Allow using a non-empty directory
  -h, --help            Show this help message
  --version             Show the forge version
`

var (
	ErrHelp    = errors.New("help requested")
	ErrVersion = errors.New("version requested")
)

// Options describes a scaffold request parsed from the command line.
type Options struct {
	Name    string
	Module  string
	HasMake bool
	Force   bool
}

// ParseArgs parses command line arguments (excluding the program name).
func ParseArgs(args []string) (Options, error) {
	var opts Options

	if len(args) == 0 {
		return opts, fmt.Errorf("missing command")
	}

	switch args[0] {
	case "new":
	case "help", "-h", "--help":
		return opts, ErrHelp
	case "version", "--version":
		return opts, ErrVersion
	default:
		return opts, fmt.Errorf("unknown command: %s", args[0])
	}

	positionals, err := parseFlags(args[1:], &opts)
	if err != nil {
		return opts, err
	}

	if len(positionals) > 1 {
		return opts, fmt.Errorf("unexpected argument: %s", positionals[1])
	}
	if opts.Name == "" && len(positionals) == 1 {
		opts.Name = positionals[0]
	}

	if opts.Name == "" {
		return opts, fmt.Errorf("missing project name")
	}
	if err := validateName(opts.Name); err != nil {
		return opts, err
	}

	return opts, nil
}

// parseFlags walks args, allowing flags to be interspersed with positionals
// and supporting both "-m value" and "--module=value" forms.
func parseFlags(args []string, opts *Options) ([]string, error) {
	var positionals []string

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positionals = append(positionals, arg)
			continue
		}

		name, value, hasValue := strings.Cut(arg, "=")

		switch name {
		case "-m", "--module":
			v, next, err := stringFlag(args, i, value, hasValue)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", name, err)
			}
			i = next
			opts.Module = v

		case "-n", "--name":
			v, next, err := stringFlag(args, i, value, hasValue)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", name, err)
			}
			i = next
			opts.Name = v

		case "-M", "--make":
			if hasValue {
				return nil, fmt.Errorf("flag %s does not take a value", name)
			}
			opts.HasMake = true

		case "-f", "--force":
			if hasValue {
				return nil, fmt.Errorf("flag %s does not take a value", name)
			}
			opts.Force = true

		case "-h", "--help":
			return nil, ErrHelp

		case "--version":
			return nil, ErrVersion

		default:
			return nil, fmt.Errorf("unknown option: %s", name)
		}
	}

	return positionals, nil
}

// stringFlag returns the value of a string flag, consuming the next argument
// when the flag uses the "-m value" form.
func stringFlag(args []string, i int, inline string, hasInline bool) (string, int, error) {
	if hasInline {
		return inline, i, nil
	}
	if i+1 >= len(args) {
		return "", i, fmt.Errorf("missing value")
	}
	return args[i+1], i + 1, nil
}

func validateName(name string) error {
	if strings.ContainsAny(name, " \t\n") {
		return fmt.Errorf("invalid project name %q: must not contain whitespace", name)
	}
	if strings.HasPrefix(name, "-") {
		return fmt.Errorf("invalid project name %q: must not start with '-'", name)
	}
	return nil
}
