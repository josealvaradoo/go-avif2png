.PHONY: build clean test test-coverage run install lint help

# Binary name
BINARY_NAME=avif2png
BUILD_DIR=bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOVET=$(GOCMD) vet

## build: Build the binary
build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/avif2png

## clean: Clean build artifacts
clean:
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

## test: Run tests
test:
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -cover -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## install: Install the binary
install: build
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)

## lint: Run go vet
lint:
	$(GOVET) ./...

## tidy: Tidy dependencies
deps:
	$(GOMOD) tidy
	$(GOMOD) download

## help: Show this help
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
