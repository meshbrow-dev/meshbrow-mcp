# Contributing to meshbrow-mcp

Thanks for your interest in contributing!

## Development

```bash
# Build
go build -o meshbrow-mcp ./cmd/mcp-server

# Test
go test ./...
go test -race ./...

# Lint
golangci-lint run
```

## Pull Requests

1. Fork and create a feature branch
2. Write tests for new functionality
3. Ensure all tests pass: `go test ./...`
4. Ensure lint passes: `golangci-lint run`
5. Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`
6. Open a PR against `main`

## Code Style

- Follow standard Go conventions
- Use `slog` for structured logging
- Error wrapping: `fmt.Errorf("doing X: %w", err)`
- Table-driven tests with `t.Run()`

## Reporting Issues

Open an issue on GitHub with:
- Steps to reproduce
- Expected vs actual behavior
- Version (`meshbrow-mcp --version`)
