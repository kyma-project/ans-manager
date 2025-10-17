# Makefile for ANS Manager CLI

.PHONY: build test clean help install

# Build the CLI
build:
	mkdir -p bin
	go build -o bin/ans-cli ./cmd/ans-cli

# Build with race detection for development
build-dev:
	mkdir -p bin
	go build -race -o bin/ans-cli ./cmd/ans-cli

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f bin/ans-cli coverage.out coverage.html

# Install the CLI to $GOPATH/bin
install:
	go install ./cmd/ans-cli

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	golangci-lint run

# Show help
help:
	@echo "Available targets:"
	@echo "  build              - Build the ans-cli binary"
	@echo "  build-dev          - Build with race detection"
	@echo "  deps               - Install/update dependencies"
	@echo "  test               - Run tests"
	@echo "  test-coverage      - Run tests with coverage report"
	@echo "  clean              - Clean build artifacts"
	@echo "  install            - Install CLI to GOPATH/bin"
	@echo "  fmt                - Format code"
	@echo "  lint               - Run linter"
	@echo "  example-user       - Example: send notification to users"
	@echo "  example-subaccount - Example: send notification to subaccounts"
	@echo "  example-globalaccount - Example: send notification to global accounts"
	@echo "  help               - Show this help"

# Example usage targets
example-user:
	./bin/ans-cli notify user admin@example.com developer@example.com \
		--region cf-eu12 \
		--subject "Application Deployed Successfully" \
		--body "Your application has been deployed to production"

example-subaccount:
	./bin/ans-cli notify subaccount 12345678-1234-1234-1234-123456789012 87654321-4321-4321-4321-210987654321 \
		--region cf-eu12 \
		--subject "Resource Alert" \
		--body "Resource usage has exceeded threshold"

example-globalaccount:
	./bin/ans-cli notify globalaccount 11111111-2222-3333-4444-555555555555 \
		--region cf-eu10-canary \
		--subject "System Maintenance" \
		--body "Scheduled maintenance will begin at 2 AM UTC"