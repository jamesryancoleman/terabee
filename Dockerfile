# syntax=docker/dockerfile:1

FROM alpine

RUN apk update
RUN apk add --no-cache git make musl-dev go

WORKDIR /terabee
COPY util/http/server.go /terabee/

RUN go build server.go