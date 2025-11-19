# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build deps
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o main cmd/app/main.go

# Final stage
FROM alpine:3.19

WORKDIR /app

# tini recommended to handle signals correctly
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/main .
COPY config ./config
COPY wait-for.sh .
COPY start.sh .

EXPOSE 8080

ENTRYPOINT ["/app/start.sh"]
