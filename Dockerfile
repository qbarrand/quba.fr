FROM golang:1.17-alpine as builder

RUN ["apk", "add", "gcc", "git", "imagemagick-dev", "make", "musl-dev", "pkgconfig", "vips-dev"]

WORKDIR /usr/src/app

COPY . .

RUN ["make"]

FROM alpine

COPY --from=builder /usr/src/app/quba-fr /quba-fr

RUN ["apk", "add", "--no-cache", "imagemagick", "vips"]

EXPOSE 8080/tcp

LABEL org.opencontainers.image.source https://github.com/qbarrand/quba.fr

ENTRYPOINT ["/quba-fr"]
