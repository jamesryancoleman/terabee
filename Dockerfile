# syntax=docker/dockerfile:1

FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get -y install golang-go docker

WORKDIR /terabee

# Copy the handler server
COPY util/http/server.go /terabee/
COPY go.* /terabee/

# Copy the terabee LXL driver image
COPY Dockerfile /terabee/

# RUN rc-update add docker boot

# Build the go server
# RUN go build server.go

# build the docker image locally
