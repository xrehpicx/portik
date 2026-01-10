
BINARY?=portik
# Version/commit -- can be set when invoking make, or computed from git
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "")
LDFLAGS?=-X 'portik/internal/cli.version=$(VERSION)' -X 'portik/internal/cli.commit=$(COMMIT)'

.PHONY: build test run clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/portik

test:
	go test ./...

run:
	go run ./cmd/portik --help

clean:
	rm -rf bin dist
