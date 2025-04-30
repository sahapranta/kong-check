.PHONY: build clean install run test fmt help release build-all

BINARY_NAME=kong-explorer
BUILD_DIR=build
MAIN_GO=main.go
VERSION=0.1.0
GO_FILES=$(shell find . -type f -name "*.go")

# Default options
help:
	@echo "Kong API Explorer"
	@echo ""
	@echo "Usage:"
	@echo "  make build    - Build the binary for current OS"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make install  - Install dependencies"
	@echo "  make run      - Run the application"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format the Go code"
	@echo "  make build-all- Build for multiple platforms"
	@echo "  make release  - Create release artifacts"

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	@echo "Dependencies installed successfully"

# Build the binary for current OS
build:
	@echo "Building Kong API Explorer..."
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_GO)
	@echo "Binary built successfully: $(BUILD_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	go clean
	@echo "Cleaned successfully"

# Run the application
run:
	@echo "Running Kong API Explorer..."
	go run $(MAIN_GO)

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Format Go code
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

# Build for multiple platforms
build-all: clean
	@echo "Building for multiple platforms..."
	mkdir -p $(BUILD_DIR)
	# Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_GO)
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_GO)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_GO)
	# Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_GO)
	@echo "Multi-platform build completed"

# Create release artifacts
release: build-all
	@echo "Creating release artifacts..."
	mkdir -p $(BUILD_DIR)/release
	cp .env.example $(BUILD_DIR)/release/.env.example
	cp README.md $(BUILD_DIR)/release/
	cd $(BUILD_DIR) && \
		tar -czf release/$(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz $(BINARY_NAME)-linux-amd64 && \
		tar -czf release/$(BINARY_NAME)-darwin-amd64-$(VERSION).tar.gz $(BINARY_NAME)-darwin-amd64 && \
		tar -czf release/$(BINARY_NAME)-darwin-arm64-$(VERSION).tar.gz $(BINARY_NAME)-darwin-arm64 && \
		zip -q release/$(BINARY_NAME)-windows-amd64-$(VERSION).zip $(BINARY_NAME)-windows-amd64.exe
	@echo "Release artifacts created in $(BUILD_DIR)/release/"

# Default
.DEFAULT_GOAL := help