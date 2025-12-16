FROM golang:1.21-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git

WORKDIR /app

# Копирование go модулей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копирование скомпилированного бинарника из builder stage
COPY --from=builder /app/main .

# Запуск приложения
CMD ["./main"]
