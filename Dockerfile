# syntax=docker/dockerfile:1
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app .

CMD sh -c "echo 'Waiting for dependencies...'; sleep 15; ./app"