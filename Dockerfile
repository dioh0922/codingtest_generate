
FROM golang:1.24

WORKDIR /app

COPY . .

#Go Modules（Goの依存管理システム）を有効
ENV GO111MODULE=on
#C言語に依存する部分を無効化
ENV CGO_ENABLED=0

EXPOSE 80