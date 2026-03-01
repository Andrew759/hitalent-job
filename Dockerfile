ARG GOLANG_VERSION='1.25.7'

FROM golang:${GOLANG_VERSION}-alpine

# Установка зависимостей и Delve
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh build-base && \
    go install github.com/go-delve/delve/cmd/dlv@master

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

# Компиляция бинарника в корень (/main), а не в текущую папку (/app/main)
RUN go build -gcflags="all=-N -l" -o /main ./cmd/main.go

EXPOSE 8080 2345

# Запуск бинарника из корня с Delve
CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/main", "--continue"]
