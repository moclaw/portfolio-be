.PHONY: build run test clean swagger docker-up docker-down docker-build docker-logs

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

# ==================== Docker Commands ====================

# Start all services with Docker Compose
docker-up:
	docker-compose up -d

# Stop all services
docker-down:
	docker-compose down

# Stop and remove volumes
docker-clean:
	docker-compose down -v

# Build Docker images
docker-build:
	docker-compose build

# View logs of all services
docker-logs:
	docker-compose logs -f

# View logs of specific service
docker-logs-backend:
	docker-compose logs -f backend

docker-logs-postgres:
	docker-compose logs -f postgres

docker-logs-redis:
	docker-compose logs -f redis

docker-logs-localstack:
	docker-compose logs -f localstack

# Restart backend service
docker-restart-backend:
	docker-compose restart backend

# Execute shell in backend container
docker-shell:
	docker-compose exec backend sh

# Check status of all services
docker-status:
	docker-compose ps

# Run database migrations in container
docker-migrate:
	docker-compose exec backend ./app migrate

# Seed database in container
docker-seed:
	docker-compose exec backend ./seed

# Rebuild and restart only backend
docker-rebuild-backend:
	docker-compose build backend
	docker-compose up -d backend

# LocalStack S3 commands
localstack-create-bucket:
	aws --endpoint-url=http://localhost:4566 s3 mb s3://portfolio-uploads

localstack-list-buckets:
	aws --endpoint-url=http://localhost:4566 s3 ls

localstack-list-files:
	aws --endpoint-url=http://localhost:4566 s3 ls s3://portfolio-uploads --recursive

# Full development setup
dev-setup: docker-up
	@echo "Waiting for services to start..."
	@sleep 10
	@echo "Development environment is ready!"
	@echo "Backend: http://localhost:8080"
	@echo "PostgreSQL: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "LocalStack S3: http://localhost:4566"

