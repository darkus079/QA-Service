FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates netcat-openbsd

RUN wget -O /usr/local/bin/goose https://github.com/pressly/goose/releases/download/v3.21.1/goose_linux_x86_64 && \
    chmod +x /usr/local/bin/goose

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/migrations ./migrations

COPY scripts/migrate.sh ./migrate.sh
RUN chmod +x ./migrate.sh

EXPOSE 8080

CMD ["./main"]
