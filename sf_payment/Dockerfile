# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy source code including vendor directory
COPY . .

# Build the application using vendor directory (offline)
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o sf-payment ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates curl

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/sf-payment .

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD curl -f http://localhost:8082/health || exit 1

# Run the application
CMD ["./sf-payment"]
