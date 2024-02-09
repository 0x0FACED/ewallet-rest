.PHONY: all
all: help test run

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: run
run:
	@echo "Running application..."
	@go run ./cmd/server/main.go

.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "----"
	@echo "Available targets:"
	@echo "  test        Run tests"
	@echo "  run         Run the application"
