# syntax=docker/dockerfile:1

FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get -y install golang-go

RUN apt-get -y install apt-transport-https ca-certificates curl gnupg2 software-properties-common
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
RUN add-apt-repository "deb [arch=arm64] https://download.docker.com/linux/debian $(lsb_release -cs) stable"

RUN apt-get update
RUN apt-get -y install docker-ce

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
