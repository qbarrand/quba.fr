FROM golang:1.16-alpine as builder

RUN ["apk", "add", "gcc", "git", "make", "musl-dev", "pkgconfig", "vips-dev"]

RUN ["mkdir", "/build"]
WORKDIR /build

COPY . .

RUN ["make"]

FROM alpine

RUN ["mkdir", "/app"]

COPY --from=builder /build/quba-fr /app
COPY data/webroot /app/webroot

RUN ["apk", "--no-cache", "add", "vips"]

EXPOSE 8080/tcp

LABEL org.opencontainers.image.source https://github.com/qbarrand/quba.fr

WORKDIR /app

ENTRYPOINT ["/app/quba-fr"]
