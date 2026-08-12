# forge

Forge is a CLI tool that scaffolds new Go projects with a clean, standard layout:

```
<app>/
├── cmd/
│   └── <app>/main.go
├── internal/
├── pkg/
├── go.mod
├── README.md
└── .gitignore
```

## Install

```bash
go install github.com/mcchukwu/forge@main
```

## Usage

```bash
forge new <app-name> [flags]
```

### Example

```bash
forge new my-app
```

## Commands

| Command | Description |
| --- | --- |
| `new <name>` | Create a new Go project |
| `help` | Show usage |
| `version` | Show the forge version |

## Flags

| Flag | Description |
| --- | --- |
| `-m, --module <name>` | Go module name (defaults to the project name) |
| `-n, --name <name>` | Override the project name |
| `-M, --make` | Also generate a Makefile |
| `-f, --force` | Allow using a non-empty directory |
| `-h, --help` | Show usage |
| `--version` | Show the forge version |

### Example with options

```bash
forge new my-app -M -m github.com/you/my-app
```

The generated project is a ready-to-run Go module, optionally initialized with
git and a Makefile providing `run`, `build`, `clean`, and `test` targets.
