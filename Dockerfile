# Сборка
FROM golang:1.23-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o todo-backend ./backend

# Финальный образ — тот же Debian (лёгкий и совместимый)
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates sqlite3 && rm -rf /var/lib/apt/lists/*
WORKDIR /root/

COPY --from=builder /app/todo-backend .
COPY --from=builder /app/frontend ./frontend

RUN mkdir -p /app/data

EXPOSE 8080

CMD ["./todo-backend"]