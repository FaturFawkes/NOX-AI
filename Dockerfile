FROM golang:1.21 AS builder
WORKDIR /app

ENV TZ="Asia/Jakarta"

COPY . .

# RUN go mod tidy

RUN go build -o noxai cmd/main.go

EXPOSE 8000

CMD ["./noxai"]