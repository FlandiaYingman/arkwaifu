# syntax=docker/dockerfile:1

# service.Dockerfile is used for build the service part of ArkWaifu. It does not contain the updateloop part, therefore
# python is not needed. The image size is only 28.79 MiB for reference (built on 2022/03/03).

FROM golang:1.24-alpine AS builder
WORKDIR /app

RUN apk update && \
    apk add --no-cache gcc libc-dev

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o service github.com/flandiayingman/arkwaifu/cmd/service


FROM alpine AS deploy

WORKDIR /app
COPY --from=builder /app/service ./service

EXPOSE 7080
WORKDIR /app
CMD [ "/app/service" ]