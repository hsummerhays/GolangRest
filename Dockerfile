# Build stage
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies (modernc.org/sqlite uses pure Go, but it doesn't hurt)
RUN apk add --no-cache gcc musl-dev

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o golangrest ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/golangrest .

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./golangrest"]
