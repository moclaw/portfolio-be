# Stage 1: Build
FROM golang:1.24-alpine AS builder

# Install only essential packages
RUN apk add --no-cache git ca-certificates && \
    apk cache clean

# Set working directory
WORKDIR /app

# Copy go mod and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Install swag for swagger docs generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate swagger docs
RUN swag init -g cmd/server/main.go -o docs

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -a -installsuffix cgo \
    -ldflags='-w -s -extldflags "-static"' \
    -o app ./cmd/server

# Stage 2: Run
FROM alpine:latest

# Install CA certificates and create user
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1001 -S appuser && \
    adduser -u 1001 -S appuser -G appuser

# Set working directory and create directories
WORKDIR /app
RUN mkdir -p /app/logs /app/data && \
    chown -R appuser:appuser /app

# Copy CA certificates and binary from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/app ./app
RUN chown appuser:appuser /app/app && chmod +x /app/app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 5303

# Run the binary
CMD ["./app"]
