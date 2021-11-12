FROM golang:1.18-alpine as go-builder

RUN ["apk", "add", "gcc", "git", "imagemagick-dev", "make", "musl-dev", "pkgconfig", "vips-dev"]

WORKDIR /usr/src/app

RUN ["mkdir", "/web-src"]

COPY cmd/ cmd/
COPY config/ config/
COPY img-src/ img-src/
COPY internal/ internal/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
COPY pkg/ pkg/

RUN ["make", "server", "img-out"]

FROM node:17-alpine as node-builder

RUN ["apk", "add", "make"]

RUN ["mkdir", "/build"]
WORKDIR /build

RUN ["mkdir", "dist"]

COPY config/ config/
COPY Makefile .
COPY package.json .
COPY package-lock.json .
COPY tsconfig.json .
COPY webpack.config.js .
COPY web-src/ web-src/

RUN ["npm", "install", "."]
RUN ["make", "webapp"]

FROM alpine

COPY --from=go-builder /usr/src/app/server /server
COPY --from=go-builder /usr/src/app/img-out /img-out
COPY --from=node-builder /build/dist /dist

EXPOSE 8080/tcp

LABEL org.opencontainers.image.source https://github.com/qbarrand/quba.fr

ENTRYPOINT ["/server", "-img-out-dir", "img-out", "-webroot-dir", "dist"]
