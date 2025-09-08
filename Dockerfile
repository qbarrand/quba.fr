FROM golang:1.25.1-alpine as go-builder

WORKDIR /usr/src/app

COPY backgrounds.mk .
COPY cmd/ cmd/
COPY img-src/ img-src/
COPY go.mod .
COPY go.sum .

RUN ["apk", "add", "gcc", "imagemagick-dev", "imagemagick-heic", "imagemagick-jpeg", "imagemagick-webp", "make", "musl-dev"]
RUN ["make", "-f", "backgrounds.mk", "backgrounds/backgrounds.json"]

FROM python:3 as python-builder

COPY fa-src/ fa-src/

RUN ["pip", "install", "fonttools[woff]"]
RUN ["make", "-C", "fa-src"]

FROM node:24-alpine as node-builder

COPY --from=go-builder /usr/src/app /build
WORKDIR /build

COPY fa-src fa-src
COPY Makefile .
COPY package.json .
COPY package-lock.json .
COPY web-src web-src
COPY webpack.config.js .

COPY --from=python-builder /fa-src/fa-brands.woff2 fa-src/
COPY --from=python-builder /fa-src/fa-solid.woff2 fa-src/

RUN ["npm", "install", "."]
RUN ["apk", "add", "make"]
RUN ["make"]

FROM scratch

COPY --from=node-builder /build/dist /
