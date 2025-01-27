# dev stage
FROM golang:1.23-bullseye AS dev

# Устанавливаем рабочую директорию
WORKDIR /app

# Устанавливаем необходимые системные зависимости
RUN apt-get update && apt-get install -y \
    protobuf-compiler \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Устанавливаем protoc-gen-go и protoc-gen-go-grpc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Экспортируем PATH для доступа к go-бинарникам
ENV PATH=$PATH:/go/bin

# Проверяем наличие protoc и плагинов
RUN protoc --version && protoc-gen-go --version && protoc-gen-go-grpc --version

# Копируем go.mod в контейнер
COPY go.mod ./

# Генерируем go.sum с помощью go mod tidy
RUN go mod tidy

# Копируем весь код
COPY . .

# Собираем приложение
# RUN go build -o ./tmp/main ./cmd

# Экспонируем порт
EXPOSE 8080

# Запускаем скомпилированное приложение
# CMD ["./tmp/main"]

# Оставляем контейнер запущенным для работы разработчика
CMD ["tail", "-f", "/dev/null"]

# test stage
FROM golang:1.23-bullseye AS test

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY --from=dev /app /app

# Компилируем приложение
RUN mkdir -p ./build && go build -o ./build/main ./cmd

CMD ["./build/main"]

# prod stage
FROM debian:bookworm-slim AS prod

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем только скомпилированный бинарник
COPY --from=test /app/main /app/build/main

# Экспонируем порт
EXPOSE 8080

# Команда для запуска приложения
CMD ["./main"]