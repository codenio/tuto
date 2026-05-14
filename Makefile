BINARY  := tuto
CMD     := ./cmd/tuto
OUT     := bin/$(BINARY)

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -ldflags "\
  -X github.com/codenio/tuto/internal/version.Version=$(VERSION) \
  -X github.com/codenio/tuto/internal/version.Commit=$(COMMIT) \
  -X github.com/codenio/tuto/internal/version.Date=$(DATE)"

.PHONY: all build install clean test lint fmt help

all: build

## build: compile binary to bin/tuto
build:
	@mkdir -p bin
	go build $(LDFLAGS) -o $(OUT) $(CMD)

## install: install binary via go install
install:
	go install $(LDFLAGS) $(CMD)

## test: run all tests
test:
	go test ./...

## lint: run golangci-lint (install: https://golangci-lint.run/usage/install/)
lint:
	golangci-lint run ./...

## fmt: format all Go source files
fmt:
	gofmt -w .
	@which goimports >/dev/null 2>&1 && goimports -w . || true

## clean: remove build artifacts
clean:
	rm -rf bin/

## help: list available targets
help:
	@grep -E '^## ' Makefile | sed 's/## /  /'
