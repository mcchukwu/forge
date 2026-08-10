package cli

import (
	"errors"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Options
		wantErr bool
	}{
		{"basic", []string{"new", "demo"}, Options{Name: "demo"}, false},
		{"module short", []string{"new", "demo", "-m", "example.com/demo"}, Options{Name: "demo", Module: "example.com/demo"}, false},
		{"module long", []string{"new", "demo", "--module", "example.com/demo"}, Options{Name: "demo", Module: "example.com/demo"}, false},
		{"module eq form", []string{"new", "demo", "--module=example.com/demo"}, Options{Name: "demo", Module: "example.com/demo"}, false},
		{"module before name", []string{"new", "-m", "example.com/demo", "demo"}, Options{Name: "demo", Module: "example.com/demo"}, false},
		{"name override", []string{"new", "demo", "-n", "other"}, Options{Name: "other"}, false},
		{"name eq form", []string{"new", "demo", "--name=other"}, Options{Name: "other"}, false},
		{"make short", []string{"new", "demo", "-M"}, Options{Name: "demo", HasMake: true}, false},
		{"make long", []string{"new", "demo", "--make"}, Options{Name: "demo", HasMake: true}, false},
		{"force short", []string{"new", "demo", "-f"}, Options{Name: "demo", Force: true}, false},
		{"force long", []string{"new", "demo", "--force"}, Options{Name: "demo", Force: true}, false},
		{"flags interspersed", []string{"new", "demo", "-M", "--force", "-m", "example.com/demo"}, Options{Name: "demo", Module: "example.com/demo", HasMake: true, Force: true}, false},
		{"all flags", []string{"new", "-n", "demo", "--module=example.com/demo", "--make", "--force"}, Options{Name: "demo", Module: "example.com/demo", HasMake: true, Force: true}, false},
		{"current dir", []string{"new", "."}, Options{Name: "."}, false},

		{"no args", []string{}, Options{}, true},
		{"unknown command", []string{"foo", "demo"}, Options{}, true},
		{"missing name", []string{"new"}, Options{}, true},
		{"module missing value", []string{"new", "demo", "-m"}, Options{}, true},
		{"unknown option", []string{"new", "demo", "--bogus"}, Options{}, true},
		{"too many positionals", []string{"new", "a", "b"}, Options{}, true},
		{"name with whitespace", []string{"new", "bad name"}, Options{}, true},
		{"name leading dash", []string{"new", "-demo"}, Options{}, true},
		{"make with value", []string{"new", "demo", "--make=1"}, Options{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseArgsSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want error
	}{
		{"help first", []string{"--help"}, ErrHelp},
		{"help short", []string{"-h"}, ErrHelp},
		{"help command", []string{"help"}, ErrHelp},
		{"help after new", []string{"new", "demo", "--help"}, ErrHelp},
		{"version long", []string{"--version"}, ErrVersion},
		{"version command", []string{"version"}, ErrVersion},
		{"version after new", []string{"new", "demo", "--version"}, ErrVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseArgs(tt.args)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v, want %v", err, tt.want)
			}
		})
	}
}
