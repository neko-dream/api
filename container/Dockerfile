FROM golang:1.24.2-alpine AS builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64

EXPOSE 3000

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -tags timetzdata -ldflags="-s -w" -trimpath -o server ./cmd/server

FROM alpine:latest

EXPOSE 3000

COPY --from=builder /app/server /server

RUN chmod +x /server

CMD ["./server"]
