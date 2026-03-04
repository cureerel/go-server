.PHONY: run build clean deps help \
        migrate-up migrate-down migrate-status migrate-create migrate-lint migrate-prod \
        lint fmt vet test prod-build

APP_NAME  = gotemplate
BUILD_DIR = build
CMD_DIR   = ./cmd/server

-include .env
export

# ── Help ──────────────────────────────────────────────────────
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ── Dependencies ──────────────────────────────────────────────
deps: ## Download and tidy dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# ── Run ───────────────────────────────────────────────────────
run: ## Run the application
	@echo "Starting server..."
	go run ./cmd/server/

# ── Build ─────────────────────────────────────────────────────
build: ## Build binary
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "Binary: $(BUILD_DIR)/$(APP_NAME)"

prod-build: ## Build production binary (Linux/amd64)
	@echo "Building production binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build -ldflags="-s -w" \
		-o $(BUILD_DIR)/$(APP_NAME) \
		$(CMD_DIR)/main.go
	@echo "Production binary: $(BUILD_DIR)/$(APP_NAME)"

# ── Clean ─────────────────────────────────────────────────────
clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)

# ── Migrations ────────────────────────────────────────────────
# tested
migrate-up:
	@echo "Applying migrations..."
	atlas migrate apply \
		--env local \
		--config file://atlas.hcl \
		--allow-dirty

migrate-down: ## Revert last migration
	@echo "Reverting last migration..."
	atlas migrate down \
		--env local \
		--config file://atlas.hcl

#tested
migrate-status: ## Show current migration status
	atlas migrate status \
		--env local \
		--config file://atlas.hcl

#tested
migrate-create: 
	@[ "$(NAME)" ] || ( echo "ERROR: NAME is required. Usage: make migrate-create NAME=add_users"; exit 1 )
	@echo "Creating migration: $(NAME)"
	atlas migrate new $(NAME) \
		--dir "file://migrations"

migrate-lint: ## Lint and validate migration files
	atlas migrate lint \
		--env local \
		--config file://atlas.hcl

migrate-prod: ## Apply migrations in production
	@echo "Applying production migrations..."
	atlas migrate apply \
		--env production \
		--config file://atlas.hcl

# ── Code Quality ──────────────────────────────────────────────
lint: ## Run linter
	golangci-lint run ./...

fmt: ## Format code
	gofmt -w .

vet: ## Run go vet
	go vet ./...

test: ## Run tests
	go test ./... -v -race -cover

db-reset: ## Drop and recreate public schema
	psql $(DATABASE_URL) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"