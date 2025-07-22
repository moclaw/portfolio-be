# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git and ca-certificates (needed for some go modules and AWS SDK)
RUN apk add --no-cache git ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS, sqlite for database, and curl for health checks
RUN apk --no-cache add ca-certificates sqlite curl

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Create data directory for SQLite
RUN mkdir -p /app/data

# Create a non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup && \
    chown -R appuser:appgroup /app

USER appuser

# Expose port
EXPOSE 5303

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:5303/health || exit 1

# Command to run
CMD ["./main"]