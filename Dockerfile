# 1. Этап сборки (Builder)
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей и качаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

RUN go mod tidy

# Собираем бинарник с названием "main"
RUN go build -o main .

# 2. Этап запуска (Runner)
FROM alpine:latest

WORKDIR /root/

# Копируем готовый бинарник из первого этапа
COPY --from=builder /app/main .

# Открываем порт
EXPOSE 8080

# Запускаем
CMD ["./main"]