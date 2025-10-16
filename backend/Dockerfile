# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Migration stage - separate image with only goose
FROM golang:1.24-alpine AS migrator

WORKDIR /migrations

# Install goose migration tool
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy migrations directory
COPY migrations ./

# Set entrypoint to goose
ENTRYPOINT ["/go/bin/goose"]

# Runtime stage - optimized image without goose (~20MB instead of 70MB)
FROM alpine:latest

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy only the binary (no goose, no migrations)
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
