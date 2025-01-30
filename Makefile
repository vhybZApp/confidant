# Project variables
APP_NAME = confidant
BUILD_DIR = bin
GO_FILES = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Define colors for output
GREEN  := \033[0;32m
CYAN   := \033[0;36m
YELLOW := \033[0;33m
RESET  := \033[0m

# Default target
.PHONY: all
all: test build

## 🛠 Build the Go application
.PHONY: build
build:
	@echo "$(CYAN)🚀 Building the application...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go
	@echo "$(GREEN)✅ Build complete! Binary is at $(BUILD_DIR)/$(APP_NAME)$(RESET)"

## 🧪 Run all tests in the internal directory
.PHONY: test
test:
	@echo "$(CYAN)🧪 Running tests...$(RESET)"
	@go test ./internal/... -v
	@echo "$(GREEN)✅ All tests passed!$(RESET)"

## 🔍 Run tests with coverage
.PHONY: coverage
coverage:
	@echo "$(CYAN)📊 Running tests with coverage...$(RESET)"
	@go test -cover ./internal/... -v
	@go test -coverprofile=coverage.out ./internal/...
	@go tool cover -html=coverage.out
	@echo "$(GREEN)✅ Coverage report generated!$(RESET)"

## 🔄 Run the application
.PHONY: run
run:
	@echo "$(CYAN)🚀 Running the application...$(RESET)"
	@go run ./cmd/main.go

## 🧹 Clean build artifacts
.PHONY: clean
clean:
	@echo "$(YELLOW)🧹 Cleaning up build artifacts...$(RESET)"
	@rm -rf $(BUILD_DIR) coverage.out
	@echo "$(GREEN)✅ Clean complete!$(RESET)"

## 🔍 Format & lint code
.PHONY: lint
lint:
	@echo "$(CYAN)🔍 Running GolangCI-Lint...$(RESET)"
	@golangci-lint run ./...
	@echo "$(GREEN)✅ Linting complete!$(RESET)"

## 🏗 Format the code properly
.PHONY: fmt
fmt:
	@echo "$(CYAN)🛠 Formatting code...$(RESET)"
	@go fmt ./...
	@echo "$(GREEN)✅ Formatting complete!$(RESET)"

## 🎯 Check for race conditions
.PHONY: race
race:
	@echo "$(CYAN)⚡ Running tests with race detector...$(RESET)"
	@go test -race ./internal/... -v
	@echo "$(GREEN)✅ Race test completed!$(RESET)"

## 📦 Install dependencies
.PHONY: deps
deps:
	@echo "$(CYAN)📦 Installing dependencies...$(RESET)"
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✅ Dependencies installed!$(RESET)"

## 🎭 Generate mock files (if using mockery)
.PHONY: mocks
mocks:
	@echo "$(CYAN)🎭 Generating mocks...$(RESET)"
	@mockery --all --output=internal/mocks --case=underscore
	@echo "$(GREEN)✅ Mocks generated!$(RESET)"

## 📖 Show available commands
.PHONY: help
help:
	@echo "$(CYAN)📖 Available commands:$(RESET)"
	@grep -E '^\.\w+|##' Makefile | sed -E 's/^\.PHONY: (.+)/  $(YELLOW)\1$(RESET)/' | sed -E 's/## (.+)/  $(GREEN)\1$(RESET)/'


