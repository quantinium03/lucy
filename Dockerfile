# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy && go mod download

COPY . .
RUN go build -o bin/lucy cmd/lucy/*

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=build /app/bin/lucy .

EXPOSE 5000

CMD ["./lucy"]
