DB_URL=postgresql://username:password@localhost:5432/forge_api?sslmode=disable
MIGRATIONS_DIR=db/migrations

## help: Show available commands
.PHONY: help
help:
	@echo "Usage:"
	@sed -n 's/^## //p' ${MAKEFILE_LIST} | column -t -s ':'

## run: Start docker services and run API
.PHONY: run
run:
	@echo "Starting docker services..."
	@docker compose up -d

	@echo "Running API..."
	@go run ./cmd/api

## build: Build binaries
.PHONY: build
build:
	@echo "Building cmd/api..."
	@go build -ldflags="-s" -o=./bin/api ./cmd/api

	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s" -o=./bin/api-linux-amd64 ./cmd/api

	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s" -o=./bin/api-darwin-amd64 ./cmd/api

	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s" -o=./bin/api-windows-amd64.exe ./cmd/api

## audit: Run formatting, vetting, and tests
.PHONY: audit
audit:
	@echo "Tidying and verifying modules..."
	@go mod tidy
	@go mod verify

	@echo "Formatting code..."
	@go fmt ./...

	@echo "Vetting code..."
	@go vet ./...

	@echo "Running tests..."
	@go test ./...

## database_up: Apply all pending migrations
.PHONY: database_up
database_up:
	@echo "Applying migrations..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

## database_down: Rollback one migration
.PHONY: database_down
database_down:
	@echo "Rolling back migration..."
	@migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

## migrate_create name=<name>: Create a new migration file pair
.PHONY: migrate_create
migrate_create:
	@test -n "$(name)" || (echo "Error: name is required. Usage: make migrate_create name=<migration_name>"; exit 1)
	@echo "Creating migration: $(name)..."
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

## swagger: Generate Swagger docs from annotations
.PHONY: swagger
swagger:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/api/main.go --output docs

## swagger_fmt: Format Swagger annotations
.PHONY: swagger_fmt
swagger_fmt:
	@echo "Formatting Swagger annotations..."
	@swag fmt