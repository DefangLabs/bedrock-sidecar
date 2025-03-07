
.PHONY: dependencies
dependencies:
	go mod tidy

# Build the application
build: dependencies $(shell find . -name "*.go")
	go build -o bin/bedrock-sidecar

# Run tests
.PHONY: test
test:
	go test -v ./...

# Run the server
run: build
	go run .

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf bin/
	go clean

# Run linter (requires golangci-lint)
.PHONY: lint
lint:
	golangci-lint run 