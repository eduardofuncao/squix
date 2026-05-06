.PHONY: test test-verbose lint fmt vet build ci clean

BINARY=squix

build:
	go build -o $(BINARY) ./cmd/squix

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
