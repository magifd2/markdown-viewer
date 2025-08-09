# Makefile for Markdown Viewer

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# Project details
BINARY_NAME=markdown-viewer
OUTPUT_DIR=bin
# Get the module path from go.mod
MODULE_PATH := $(shell go list -m)

# Versioning
# Get the version from the latest git tag
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
# Get the git commit hash
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")
# LDFLAGS to embed version and commit info
LDFLAGS=-ldflags "-s -w -X '$(MODULE_PATH).version=$(VERSION) (commit: $(COMMIT_HASH))'"

.PHONY: all build clean test cross-compile package-all vulncheck help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all           : Build binaries for all target platforms."
	@echo "  build         : Build the binary for the current OS and architecture."
	@echo "  run           : Build and run the application."
	@echo "  test          : Run all tests."
	@echo "  vulncheck     : Run vulnerability check."
	@echo "  clean         : Clean up build artifacts."
	@echo "  cross-compile : Cross-compile for all target platforms (macOS, Linux, Windows)."
	@echo "  package-all   : Package all cross-compiled binaries into archives."
	@echo "  help          : Display this help message."

all: cross-compile package-all

# Build for the current OS/Arch
build:
	@echo "Building for $(shell go env GOOS)/$(shell go env GOARCH)..."
	@$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Run the application
run: build
	@./$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@$(GOTEST) -v ./...

# Run vulnerability check
vulncheck:
	@echo "Running vulnerability check..."
	@$(GOCMD) run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(OUTPUT_DIR)

# Cross-compile for all target platforms
cross-compile: build-mac-universal build-linux build-windows
	@echo "Cross-compilation finished. Binaries are in the $(OUTPUT_DIR)/ directory."

# Build for Linux (amd64 & arm64)
build-linux:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(OUTPUT_DIR)/linux-amd64
	@GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/linux-amd64/$(BINARY_NAME) .
	@echo "Building for Linux (arm64)..."
	@mkdir -p $(OUTPUT_DIR)/linux-arm64
	@GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/linux-arm64/$(BINARY_NAME) .

# Build for Windows (amd64)
build-windows:
	@echo "Building for Windows (amd64)..."
	@mkdir -p $(OUTPUT_DIR)/windows-amd64
	@GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/windows-amd64/$(BINARY_NAME).exe .

# Build macOS Universal Binary
build-mac-universal:
	@echo "Building for macOS (Universal)..."
	@mkdir -p $(OUTPUT_DIR)/darwin-universal
	@GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@lipo -create -output $(OUTPUT_DIR)/darwin-universal/$(BINARY_NAME) $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64
	@codesign -s - $(OUTPUT_DIR)/darwin-universal/$(BINARY_NAME)
	@rm $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-amd64 $(OUTPUT_DIR)/$(BINARY_NAME)-darwin-arm64
	@echo "Created Universal binary at $(OUTPUT_DIR)/darwin-universal/$(BINARY_NAME)"

# Package all binaries into archives
package-all: 
	@echo "Packaging all binaries..."
	$(MAKE) package-darwin
	$(MAKE) package-linux
	$(MAKE) package-windows

# Package macOS binary
package-darwin:
	@echo "Packaging macOS binary..."
	@cd $(OUTPUT_DIR)/darwin-universal && tar -czvf ../$(BINARY_NAME)-$(VERSION)-darwin-universal.tar.gz $(BINARY_NAME)

# Package Linux binaries
package-linux:
	@echo "Packaging Linux binaries..."
	@cd $(OUTPUT_DIR)/linux-amd64 && tar -czvf ../$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz $(BINARY_NAME)
	@cd $(OUTPUT_DIR)/linux-arm64 && tar -czvf ../$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz $(BINARY_NAME)

# Package Windows binary
package-windows:
	@echo "Packaging Windows binary..."
	@cd $(OUTPUT_DIR)/windows-amd64 && zip -r ../$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BINARY_NAME).exe
