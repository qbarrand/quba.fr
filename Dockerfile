FROM golang:1.15-alpine as builder

RUN ["apk", "add", "gcc", "git", "make", "musl-dev", "pkgconfig", "vips-dev"]

RUN ["mkdir", "/build"]
WORKDIR /build

COPY ./cmd cmd
COPY ./internal internal
COPY ./pkg pkg
COPY go.* .
COPY Makefile .

RUN ["make"]

FROM alpine

RUN ["mkdir", "/app"]

COPY --from=builder /build/quba-fr /app
COPY ./webroot /app/webroot

RUN ["apk", "add", "vips"]

EXPOSE 8080/tcp

LABEL org.opencontainers.image.source https://github.com/qbarrand/quba.fr

ENTRYPOINT ["/app/quba-fr"]
