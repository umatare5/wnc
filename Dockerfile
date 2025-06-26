# Dockerfile for wnc CLI
# Multi-stage build for optimized container size

# Build stage
FROM golang:1.24-alpine AS builder

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o wnc ./cmd/main.go

# Final stage - minimal container
FROM scratch

# Copy ca-certificates for HTTPS requests to WNC controllers
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /app/wnc /wnc

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/wnc"]

# Default command shows help
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="wnc"
LABEL org.opencontainers.image.description="CLI tool for managing Cisco C9800 Wireless Network Controllers"
LABEL org.opencontainers.image.vendor="umatare5"
LABEL org.opencontainers.image.source="https://github.com/umatare5/wnc"
LABEL org.opencontainers.image.documentation="https://github.com/umatare5/wnc/blob/main/README.md"
