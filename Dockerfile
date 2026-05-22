FROM golang:1.23.8-alpine3.20 AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates
COPY go.mod ./
RUN go mod download

COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/api/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/app ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates postgresql-client bash
WORKDIR /app

COPY --from=builder /out/app /usr/local/bin/app
COPY migrations ./migrations
COPY web ./web
COPY docker-entrypoint.sh /docker-entrypoint.sh

RUN chmod +x /docker-entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]
