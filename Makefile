# Makefile for River project

# Variables
BINARY=river
MAIN_DIR=.

# Build the project
build:
	go build -o ${BINARY} ${MAIN_DIR}

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run linter
lint:
	golangci-lint run

# Fix formatting
fmt:
	gofmt -w .
	goimports -w .

# Clean build files
clean:
	rm -f ${BINARY}

# Install dependencies
deps:
	go mod tidy

# Run the application
run: build
	./${BINARY}

.PHONY: build test test-coverage lint fmt clean deps run