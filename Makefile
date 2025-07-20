.PHONY: build run test clean swagger

# Build the application
build:
	go build -o bin/portfolio-be cmd/server/main.go

# Run the application
run: swagger
	go run cmd/server/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Generate swagger documentation
swagger:
	export PATH=$$PATH:$$(go env GOPATH)/bin && swag init -g cmd/server/main.go -o ./docs

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install swag if not already installed
install-swag:
	go install github.com/swaggo/swag/cmd/swag@latest

# Development server with live reload (requires air)
dev: swagger
	$$(go env GOPATH)/bin/air

# Install air for live reload
install-air:
	go install github.com/air-verse/air@latest

# Seed database with sample data
seed:
	go run cmd/seed/main.go
