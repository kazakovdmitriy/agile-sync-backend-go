# Первый этап: сборка
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Установка bash (для скриптов) и ca-certificates
RUN apk add --no-cache bash ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем и сервер, и миграционный CLI
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/agile_sync
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./cmd/migrate

# Второй этап: создание минимального образа для запуска
FROM alpine:latest

RUN apk --no-cache add ca-certificates bash

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations/

COPY --from=builder /app/scripts/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]