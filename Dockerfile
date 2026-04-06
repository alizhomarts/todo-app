# Используем базовый образ golang для сборки приложения
FROM golang:1.26-alpine AS build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /todo-app

# Устанавливаем git для загрузки зависимостей
RUN apk add --no-cache git

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем все исходные файлы приложения
COPY . .

# Собираем приложение
WORKDIR /todo-app/cmd/app
RUN go build -o main .

# Финальный образ на основе Alpine для уменьшения размера
FROM alpine:latest

# Устанавливаем рабочую директорию
WORKDIR /todo-app

RUN apk add --no-cache ca-certificates

# Копируем исполняемый файл приложения
COPY --from=build /todo-app/cmd/app/main /todo-app/main

# Копируем миграции базы данных
COPY --from=build /todo-app/database/migrations /todo-app/database/migrations

# Устанавливаем команду по умолчанию для запуска приложения
CMD ["/todo-app/main"]