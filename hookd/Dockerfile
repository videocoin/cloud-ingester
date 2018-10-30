FROM golang:1.11-rc-alpine AS builder

RUN apk update && apk add --update build-base alpine-sdk musl-dev musl

WORKDIR /go/src/gitlab.videocoin.io/videocoin/ingester/hookd

ADD . ./

ENV GO111MODULE off

RUN make build-alpine

FROM alpine:latest AS release

COPY --from=builder /go/src/gitlab.videocoin.io/videocoin/ingester/hookd/bin/hookd ./

ENTRYPOINT ./hookd

