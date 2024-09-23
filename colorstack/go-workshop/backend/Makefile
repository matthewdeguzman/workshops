.PHONY: all build run clean

all: build

build:
	@echo "Building..."
	@go build -o bin/ ./cmd/...

run-build:
	@echo "Running build..."
	@./bin/app

run:
	@echo "Running..."
	@go run ./cmd/...

clean:
	@echo "Cleaning..."
	@rm -rf bin/
