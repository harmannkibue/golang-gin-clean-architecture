# 🏗 Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies for Go modules
RUN apk add --no-cache git

# Copy go.mod and go.sum separately to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source files
COPY . .

# Build binary with optimizations (strip debug info & enable optimizations)
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/app/main.go

# 🚀 Runtime stage (small final image)
FROM alpine:latest

WORKDIR /app

# Install any required runtime dependencies
RUN apk add --no-cache bash

# Copy only the necessary files from the builder stage
COPY --from=builder /app/main .
COPY config ./config
COPY wait-for.sh .
COPY start.sh .
COPY ./migrations ./migrations

# Make scripts executable
RUN chmod +x /app/start.sh /app/wait-for.sh

# Expose port 8080
EXPOSE 8080

# Define the startup command
ENTRYPOINT ["/app/start.sh"]

