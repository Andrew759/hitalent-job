ARG GOLANG_VERSION='1.25.7'

FROM golang:${GOLANG_VERSION}-alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./
COPY ./.env.example .env

RUN go build -o main ./cmd/main.go

# Создается финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируется бинарник из стадии сборки
COPY --from=builder . .

EXPOSE 8080

# Запуск приложения
CMD ["./app/main"]
