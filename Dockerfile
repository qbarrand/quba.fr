FROM golang:1.15-alpine as builder

RUN ["apk", "add", "git", "make"]

RUN ["mkdir", "/build"]
WORKDIR /build

COPY ./cmd cmd
COPY ./internal internal
COPY ./pkg pkg
COPY go.* .
COPY Makefile .

RUN ["make"]

FROM alpine

COPY --from=builder /build/quba-fr /

ENTRYPOINT ["/quba-fr"]
