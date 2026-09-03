.PHONY: test test-verbose lint fmt vet build build-lite build-minimal ci clean

BINARY=squix
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS=-ldflags "-s -w -X main.Version=$(VERSION)"

LITE_TAGS=noduckdb,nosnowflake,noracle
MINIMAL_TAGS=noduckdb,nosnowflake,noracle,noclickhouse,nofirebird,nosqlserver

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/squix

build-lite:
	go build $(LDFLAGS) -tags "$(LITE_TAGS)" -o $(BINARY)-lite ./cmd/squix

build-minimal:
	go build $(LDFLAGS) -tags "$(MINIMAL_TAGS)" -o $(BINARY)-minimal ./cmd/squix

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
	rm -f $(BINARY) $(BINARY)-lite $(BINARY)-minimal
