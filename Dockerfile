# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main cmd/app/main.go
RUN apk add curl

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY config ./config
COPY wait-for.sh .
COPY start.sh .

EXPOSE 8080
CMD ["/app/main"]

ENTRYPOINT ["/app/start.sh"]
