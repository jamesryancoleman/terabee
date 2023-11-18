# syntax=docker/dockerfile:1

FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get -y install golang-go

# Add Docker's official GPG key:
RUN apt-get update
RUN apt-get install -y ca-certificates curl gnupg
RUN install -m 0755 -d /etc/apt/keyrings
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
RUN chmod a+r /etc/apt/keyrings/docker.gpg

# Add the repository to Apt sources:
RUN echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update
 
RUN apt-get install -y docker-ce docker-ce-cli containerd.io
 # docker-buildx-plugin docker-compose-plugin

WORKDIR /terabee

# Copy the handler server
COPY util/http/server.go /terabee/
COPY go.* /terabee/

# Copy the terabee LXL driver image
COPY Dockerfile /terabee/

# Build the go server
# RUN go build server.go

# build the docker image locally
# RUN docker build -t terabee .
