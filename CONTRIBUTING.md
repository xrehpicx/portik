# Contributing

Thanks for contributing to **portik**.

## Requirements
- Go 1.22+
- On Linux: `ss` and `ps` in PATH
- On macOS: `lsof` and `ps` in PATH
- Optional: `docker` in PATH for `--docker` features

## Local dev
```bash
go test ./...
go run ./cmd/portik --help