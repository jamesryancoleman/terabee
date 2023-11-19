# syntax=docker/dockerfile:1

FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get -y install golang-go

WORKDIR /terabee

COPY driver/convert.go /terabee/driver/