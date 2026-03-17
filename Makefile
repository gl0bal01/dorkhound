VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINARY := dorkhound

.PHONY: build install release clean test lint fmt

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/dorkhound

install:
	go install $(LDFLAGS) ./cmd/dorkhound

test:
	go test ./... -v

release:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-linux-amd64 ./cmd/dorkhound
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-darwin-amd64 ./cmd/dorkhound
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY)-darwin-arm64 ./cmd/dorkhound
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-windows-amd64.exe ./cmd/dorkhound

clean:
	rm -f $(BINARY) $(BINARY)-*

lint:
	go vet ./...
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Files not formatted:"; echo "$$unformatted"; exit 1; \
	fi

fmt:
	gofmt -w .
