
FROM golang:1.24

WORKDIR /app

#Go Modules（Goの依存管理システム）を有効
ENV GO111MODULE=on
#C言語に依存する部分を無効化
ENV CGO_ENABLED=0

COPY . .
RUN go mod tidy

RUN go build -o /app/main .

EXPOSE 80