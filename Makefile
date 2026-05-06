.PHONY: test test-verbose lint fmt vet build ci clean

BINARY=squix
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-X main.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/squix

test:
	go test ./...

test-verbose:
	go test -v ./...

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

vet:
	go vet ./...

ci: fmt vet lint test

clean:
	rm -f $(BINARY)
