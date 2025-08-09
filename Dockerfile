# Многоэтапная сборка для минимизации размера образа
FROM golang:1.24.4-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git ca-certificates

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кеширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.version=1.0.0" -o criage .

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk add --no-cache ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -g 1001 -S criage && \
    adduser -u 1001 -S criage -G criage

# Создаем рабочую директорию
WORKDIR /home/criage

# Копируем исполняемый файл из стадии сборки
COPY --from=builder /app/criage /usr/local/bin/criage

# Создаем директории для конфигурации и данных
RUN mkdir -p /home/criage/.criage && \
    chown -R criage:criage /home/criage

# Переключаемся на непривилегированного пользователя
USER criage

# Переменные окружения
ENV CRIAGE_VERSION=1.0.0
ENV CRIAGE_HOME=/home/criage/.criage

# Открываем порт (если потребуется для будущих версий)
EXPOSE 8080

# Точка входа
ENTRYPOINT ["criage"]

# Команда по умолчанию
CMD ["--help"]

# Метаданные образа
LABEL maintainer="Criage Team"
LABEL version="1.0.0"
LABEL description="Высокопроизводительный пакетный менеджер Criage"
LABEL org.opencontainers.image.source="https://github.com/criage-oss/criage-client"
LABEL org.opencontainers.image.documentation="https://criage.ru/docs.html"
LABEL org.opencontainers.image.licenses="MIT"
