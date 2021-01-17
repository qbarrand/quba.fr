FROM golang:1.15-alpine as builder

RUN ["apk", "add", "gcc", "git", "make", "musl-dev", "pkgconfig", "vips-dev"]

RUN ["mkdir", "/build"]
WORKDIR /build

COPY . /build

RUN ["make"]

FROM alpine

RUN ["mkdir", "/app"]

COPY --from=builder /build/quba-fr /app
COPY ./webroot /app/webroot
COPY ./templates /app/templates

RUN ["apk", "add", "vips"]

EXPOSE 8080/tcp

LABEL org.opencontainers.image.source https://github.com/qbarrand/quba.fr

WORKDIR /app

ENTRYPOINT ["/app/quba-fr"]
