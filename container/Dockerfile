FROM golang:1.23.1-alpine3.20 as base

ARG API_NAME
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=arm64

EXPOSE 3000

WORKDIR /app

COPY ./go.mod /app/go.mod
RUN go mod download

COPY . /app
RUN go build -tags timetzdata -ldflags="-s -w" -trimpath -o server .

FROM --platform=linux/arm64 alpine:latest
ARG ENV ${ENV:-dev}
COPY --from=base /app/server /app/server
CMD /app/server
