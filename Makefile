# Define variables
GO := go
LINTER := golangci-lint
PKGS := $(shell go list ./...)

# Load .db.env variables
ifeq (,$(wildcard .env))
$(error ".env file not found")
endif
include .env
export $(shell sed 's/=.*//' .env)

POSTGRESQL_URL="postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST_NAME):$(POSTGRES_HOST_PORT)/$(POSTGRES_DB)?sslmode=disable"

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

sqlc-generate:
	@sqlc generate

new-migration:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)

migrateup:
	migrate -path internal/db/migrations -database ${POSTGRESQL_URL} -verbose up

migratedown:
	migrate -path internal/db/migrations -database ${POSTGRESQL_URL} -verbose down

air: docs
	@air

docs:
	@rm -rf docs/
	@swag init --pd
