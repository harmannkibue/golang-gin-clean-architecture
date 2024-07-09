# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o main cmd/app/main.go
#RUN apk add curl
#RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
#COPY --from=builder  /app/migrate.linux-amd64 ./migrate
COPY config ./config
COPY wait-for.sh .
COPY start.sh .
COPY ./migrations ./migrations

EXPOSE 8082
CMD ["/app/main"]

ENTRYPOINT ["/app/start.sh"]
