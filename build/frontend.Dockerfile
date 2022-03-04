# syntax=docker/dockerfile:1

# service.Dockerfile is used for build the service part of ArkWaifu. It does not contain the updateloop part, therefore
# python is not needed. The image size is only 28.79 MiB for reference (built on 2022/03/03).

FROM node:17-alpine AS builder
WORKDIR /app

COPY ./web/arkwaifu/package.json ./web/arkwaifu/package-lock.json ./
RUN npm install

COPY ./web/arkwaifu ./
RUN npm run build


FROM nginx:alpine AS deploy

COPY --from=builder /app/dist /usr/share/nginx/html
COPY ./build/frontend.nginx.conf /etc/nginx/nginx.conf

EXPOSE 80
