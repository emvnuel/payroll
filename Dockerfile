ARG GO_VERSION=1.21
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*
RUN mkdir -p /api
WORKDIR /api
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o ./app ./main.go

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN mkdir -p /api
RUN mkdir -p /api/resources
WORKDIR /api
COPY --from=builder /api/app .
ENTRYPOINT ["./app"]