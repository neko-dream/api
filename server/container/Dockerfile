FROM oven/bun:1.2.12-alpine AS admin-ui-builder

ARG ENV=dev
ARG API_URL=https://api-dev.kotohiro.com

WORKDIR /admin-ui
COPY admin-ui/package*.json admin-ui/bun.lockb ./
RUN bun install
COPY admin-ui ./
RUN VITE_API_URL=${API_URL} bun build:dev

FROM golang:1.25.1-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

EXPOSE 3000

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=admin-ui-builder /static/admin-ui /static/admin-ui

RUN go build -tags timetzdata -ldflags="-s -w" -trimpath -o server ./cmd/server

FROM alpine:latest

EXPOSE 3000

COPY --from=builder /app/server /server
RUN chmod +x /server

CMD ["./server"]
