# syntax=docker/dockerfile:1

# updateloop.Dockerfile is used for build the updateloop part of ArkWaifu. It contains the updateloop part, therefore
# python is needed.

FROM golang:1.20-bullseye AS builder
WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y gcc libc-dev
COPY ./cmd  ./cmd
COPY ./internal ./internal
RUN go build -o updateloop github.com/flandiayingman/arkwaifu/cmd/updateloop

FROM python:3.10-slim-bullseye AS deploy

WORKDIR /app
COPY ./tools/extractor ./tools/extractor

WORKDIR /app/tools/extractor
RUN pip install --no-cache-dir -r requirements.txt

WORKDIR /app
COPY --from=builder /app/updateloop ./updateloop
CMD [ "/app/updateloop" ]