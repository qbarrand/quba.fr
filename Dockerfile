FROM golang:1.22.5-alpine as go-builder

WORKDIR /usr/src/app

COPY cmd/ cmd/
COPY config/ config/
COPY img-src/ img-src/
COPY internal/ internal/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
COPY pkg/ pkg/

RUN ["apk", "add", "gcc", "git", "imagemagick-dev", "make", "musl-dev", "pkgconfig", "vips-dev"]
RUN ["make", "server", "img-out"]

FROM python:3 as python-builder

COPY fa-src/ fa-src/

RUN ["pip", "install", "fonttools[woff]"]
RUN ["make", "-C", "fa-src"]

FROM node:22-alpine as node-builder

RUN ["mkdir", "/build"]
WORKDIR /build

COPY config/ config/
COPY fa-src/ fa-src/
COPY Makefile .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .
COPY webpack.config.js .
COPY web-src/ web-src/

RUN ["mkdir", "dist"]

COPY --from=python-builder /fa-src/fa-brands.woff2 fa-src/
COPY --from=python-builder /fa-src/fa-solid.woff2 fa-src/

RUN ["npm", "install", "."]

RUN ["apk", "add", "make"]
RUN ["make", "webapp"]

FROM alpine

COPY --from=go-builder /usr/src/app/server /server
COPY --from=go-builder /usr/src/app/img-out /img-out
COPY --from=node-builder /build/dist /dist

EXPOSE 8080/tcp

ENTRYPOINT ["/server", "-img-out-dir", "img-out", "-webroot-dir", "dist"]
