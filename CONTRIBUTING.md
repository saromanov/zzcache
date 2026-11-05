# Contributing

Thank you for considering a contribution.

## Development

- Go 1.21 or newer.
- Run tests with the race detector: `go test -race -count=1 ./...`
- Run the linter: `golangci-lint run` (config in `.golangci.yml`).
- Format code with `gofmt`/`goimports` before committing.

## Guidelines

- Keep the public API small and additive. Breaking changes need a new major/minor version discussion first.
- Every new feature or bug fix must come with tests covering both store types (default map and radix).
- Benchmarks for performance-sensitive changes are welcome (`bench_test.go`).
- One change per pull request.

## Reporting bugs

Open a GitHub issue with the bug report template: Go version, minimal reproduction, expected vs actual behavior.