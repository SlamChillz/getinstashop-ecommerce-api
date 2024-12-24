# Define variables
GO := go
LINTER := golangci-lint
PKGS := $(shell go list ./...)

# Targets
.PHONY: all build run lint test clean

all: build

build:
	$(GO) build -o bin/server .

run: build
	./bin/server

lint: check-linter
	$(LINTER) run ./...

test:
	$(GO) test -v ./...

clean:
	rm -rf bin/

# Ensure the linter is installed
check-linter:
	@command -v $(LINTER) >/dev/null 2>&1 || { \
		echo "Installing $(LINTER)..."; \
		$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	}
