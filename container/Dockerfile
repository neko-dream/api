FROM golang:1.23.1-alpine3.20 as base

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64

EXPOSE 3000

WORKDIR /app

COPY ./go.* /app
RUN go mod download

COPY . /app
RUN go build -tags timetzdata -ldflags="-s -w" -trimpath -o server .

FROM --platform=linux/arm64 alpine:latest

COPY --from=base /app/server /app/server
CMD /app/server
