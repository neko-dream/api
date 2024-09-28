FROM golang:1.23.1-alpine3.20 as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64

EXPOSE 3000

WORKDIR /app

COPY ./go.* .
RUN go mod download

COPY . .

RUN go build -tags timetzdata -ldflags="-s -w" -trimpath -o server ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server ./server
CMD /app/server
