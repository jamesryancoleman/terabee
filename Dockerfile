# syntax=docker/dockerfile:1

FROM alpine

RUN apk update
RUN apk add --no-cache git make musl-dev go

WORKDIR /terabee

# Copy the handler server
COPY util/http/server.go /terabee/
# Copy the terabee LXL driver image
COPY Dockerfile /terabee/

# Build the go server
RUN go build server.go

# build the docker image locally
