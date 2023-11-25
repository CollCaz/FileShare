
# Simple Makefile for a Go project

# Build the application
all: build buildWindows

# Build for linux
build:
	@echo "Building..."
	@go build -o bin/FileShareServer cmd/app/main.go

# Build for windows
buildWindows:
	@echo "Bulding for windows..."
	@GOOS=windows GOARCH=amd64 go build -o bin/FileShareServer.exe cmd/app/main.go
# Run the application
run:
	@go run cmd/app/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

.PHONY: all build run test clean
		
