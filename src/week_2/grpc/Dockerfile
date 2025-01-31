FROM golang:1.23.4-alpine AS builder

COPY .. github.com/xdevspo/go-microservices/week_2/grpc/source/

# Устанавливаем рабочую директорию
WORKDIR /github.com/xdevspo/go-microservices/week_2/grpc/source/

# Копируем файлы go.mod и go.sum перед загрузкой зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY .. .

# Собираем бинарник
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

# Используем минимальный образ для продакшена
FROM alpine:latest

WORKDIR /root/

# Копируем скомпилированный бинарник
COPY --from=builder /github.com/xdevspo/go-microservices/week_2/grpc/source/bin/crud_server .

# Делаем его исполняемым
RUN chmod +x crud_server

# Запускаем сервер
CMD ["./crud_server"]