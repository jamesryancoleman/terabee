# syntax=docker/dockerfile:1

FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get -y install golang-go

WORKDIR /terabee/driver/convert

COPY driver/convert/convert_flow.go /terabee/driver/convert/