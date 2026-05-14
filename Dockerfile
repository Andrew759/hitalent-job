ARG GOLANG_VERSION='1.26.3'

FROM golang:${GOLANG_VERSION}-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh build-base && \
    go install github.com/go-delve/delve/cmd/dlv@master

WORKDIR /app

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

COPY main /main
RUN chmod +x /main

EXPOSE 8080 2345

# Запуск бинарника из корня с Delve
CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/main", "--continue"]