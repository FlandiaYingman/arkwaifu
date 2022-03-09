# syntax=docker/dockerfile:1

# updateloop.Dockerfile is used for build the updateloop part of ArkWaifu. It contains the updateloop part, therefore
# python is needed. The image size is only 28.79 MiB for reference (built on 2022/03/03).

FROM golang:1.17-alpine AS builder
WORKDIR /app

RUN apk update && \
    apk add --no-cache gcc libc-dev

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN go build -o updateloop github.com/flandiayingman/arkwaifu/cmd/updateloop


FROM python:3.10-alpine AS deploy

WORKDIR /app
COPY --from=builder /app/updateloop ./updateloop
COPY --from=builder /app/tools/extractor ./tools/extractor

WORKDIR /app/tools/extractor
RUN apk update && \
    apk add --no-cache gcc g++ libc-dev zlib-dev jpeg-dev && \
    apk add --no-cache libstdc++ libjpeg
RUN pip install --no-cache-dir pipenv && \
	pipenv install && \
	pipenv --clear
RUN apk del gcc g++ libc-dev zlib-dev jpeg-dev

WORKDIR /app
CMD [ "/app/updateloop" ]