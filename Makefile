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

completion-install: build
	@shell=$$(basename "$$SHELL"); \
	case "$$shell" in \
	  bash) ./$(BINARY) completion bash > /etc/bash_completion.d/$(BINARY) && echo "bash completion installed to /etc/bash_completion.d/$(BINARY)";; \
	  zsh)  ./$(BINARY) completion zsh > "$${fpath[1]}/_$(BINARY)" && echo "zsh completion installed";; \
	  fish) ./$(BINARY) completion fish > ~/.config/fish/completions/$(BINARY).fish && echo "fish completion installed";; \
	  *)    echo "Unsupported shell: $$shell. Run: dorkhound completion <bash|zsh|fish|powershell>";; \
	esac

release-local:
	mkdir -p dist
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/$(BINARY)-linux-amd64   ./cmd/dorkhound
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/$(BINARY)-linux-arm64   ./cmd/dorkhound
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/$(BINARY)-darwin-amd64  ./cmd/dorkhound
	GOOS=darwin  GOARCH=arm64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/$(BINARY)-darwin-arm64  ./cmd/dorkhound
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath $(LDFLAGS) -o dist/$(BINARY)-windows-amd64.exe ./cmd/dorkhound
	cd dist && for f in *; do tar czf "$$f.tar.gz" "$$f" && rm "$$f"; done
	@echo "Local release artifacts in dist/"
