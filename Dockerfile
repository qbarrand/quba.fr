FROM golang:1.16-alpine as builder

RUN ["apk", "add", "gcc", "git", "make", "musl-dev", "pkgconfig", "vips-dev"]

RUN ["mkdir", "/build"]
WORKDIR /build

COPY . .

RUN ["make"]

FROM alpine

COPY --from=builder /build/quba-fr /quba-fr

RUN ["apk", "add", "--no-cache", "vips"]

EXPOSE 8080/tcp

LABEL org.opencontainers.image.source https://github.com/qbarrand/quba.fr

ENTRYPOINT ["/quba-fr"]
