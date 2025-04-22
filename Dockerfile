FROM golang:1.24.2-alpine as go-builder

WORKDIR /usr/src/app

COPY cmd/ cmd/
COPY img-src/ img-src/
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

RUN ["apk", "add", "gcc", "git", "imagemagick-dev", "make", "musl-dev"]
RUN ["make", "backgrounds"]

FROM python:3 as python-builder

COPY fa-src/ fa-src/

RUN ["pip", "install", "fonttools[woff]"]
RUN ["make", "-C", "fa-src"]

FROM node:23-alpine as node-builder

RUN ["mkdir", "/build"]
WORKDIR /build

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

FROM scratch

COPY --from=node-builder /build/dist /
COPY --from=go-builder /usr/src/app/backgrounds /backgrounds
