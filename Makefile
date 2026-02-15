.PHONY: run build clean deps help

# Variables
APP_NAME=gotemplate
BUILD_DIR=build
CMD_DIR=./cmd/server

help: ## Show this help message
    @echo 'Usage: make [target]'
	@echo ''
    @echo 'Available targets:'
    @awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

run: ## Run the application
	@echo "Running application..."
	go run $(CMD_DIR)/main.go

build: ## Build the application binary
	@echo "Building application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go

clean: ## Clean build files
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

migrate-up: ## Run database migrations (Placeholder)
	@echo "Running migrations..."