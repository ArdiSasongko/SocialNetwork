# Build stage
FROM golang:1.23.2-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod tidy

# Copy application files
COPY . ./

# Build the application
RUN go build -o bin/main cmd/api/*.go

# Final stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/bin/main ./bin/main

# Copy necessary files
COPY --from=builder /app/.env .env

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./bin/main"]