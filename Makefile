.PHONY: all clean fmt lint test test-cover tidy check bump install examples help

EXAMPLES := $(patsubst _examples/%/,%,$(wildcard _examples/*/))
DIST_DIR := dist

# Go commands
GO := go
GOFMT := gofmt
GOFILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: fmt lint test

## fmt: Format Go source files
fmt:
	$(GOFMT) -w $(GOFILES)

## lint: Run linter
lint:
	@golangci-lint run || echo "golangci-lint not installed, skipping"

## test: Run tests
test:
	$(GO) test -v ./...

## test-cover: Run tests with coverage
test-cover:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

## tidy: Tidy go modules
tidy:
	$(GO) mod tidy

## examples: Build all examples into dist/
examples:
	@mkdir -p $(DIST_DIR)
	@for ex in $(EXAMPLES); do \
		echo "Building $$ex..."; \
		$(GO) build -o $(DIST_DIR)/$$ex ./_examples/$$ex; \
	done
	@echo "All examples built in $(DIST_DIR)/"

## clean: Remove generated artifacts
clean:
	rm -f coverage.out coverage.html
	rm -rf $(DIST_DIR)
	$(GO) clean

## check: Run fmt, lint, and test
check: fmt lint test

## bump: Tag and push next version (requires clean tree and svu)
bump:
	@$(MAKE) check
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "Working directory not clean. Commit or stash first."; \
		exit 1; \
	fi
	@version=$$(svu next); \
		git tag -a $$version -m "Version $$version"; \
		echo "Tagged $$version"; \
		git push origin $$version

## install: Install the scaffold command
install:
	$(GO) install ./cmd/scaffold

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
