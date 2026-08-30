# Agent Guidelines

## Build & Test

- Build all: `./scripts/build.sh [output-dir]`
- Build single: `go build -o bin/<name> ./cmd/<name>`
- Run e2e tests: `go test ./e2e/...`
- Check code with `go fmt` and `go vet`
- **No unit tests** — e2e only. Tests build the real binary and mock external dependencies (see `e2e/` for the pattern).

## Project Conventions

- One binary per directory under `cmd/`. Directory name = binary name.
- Update the README scripts table when adding a new CLI tool.
- `internal/` for private shared packages, `pkg/` for public/reusable.
- No external Go dependencies — stick to the standard library.

## Code Style

- Idiomatic Go (https://go.dev/doc/effective_go).
- Prefer `flag` or `os.Args` for CLI arguments.
- Return errors up the call stack, handle I/O and `os.Exit(1)` only in `main()`.
