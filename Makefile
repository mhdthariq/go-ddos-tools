.PHONY: build clean test run help

# Build the project
build:
	@echo "Building DDoS-Tools..."
	go build -o ddos-tools main.go

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -o bin/ddos-tools-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/ddos-tools-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/ddos-tools-windows-amd64.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/ddos-tools-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/ddos-tools-darwin-arm64 main.go
	@echo "Build complete! Binaries are in bin/"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f ddos-tools
	rm -rf bin/
	go clean

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run the application
run:
	@echo "Running DDoS-Tools..."
	go run main.go

# Show help
help:
	@echo "DDoS-Tools - Go Implementation"
	@echo ""
	@echo "Available targets:"
	@echo "  build      - Build the project for current platform"
	@echo "  build-all  - Build for all platforms (Linux, Windows, macOS)"
	@echo "  clean      - Remove build artifacts"
	@echo "  test       - Run tests"
	@echo "  deps       - Install/update dependencies"
	@echo "  run        - Run the application"
	@echo "  help       - Show this help message"
