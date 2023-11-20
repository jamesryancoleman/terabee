# syntax=docker/dockerfile:1

FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get -y install golang-go

WORKDIR /terabee/driver/convert

COPY convert/convert_flow.go /terabee/driver/convert/

WORKDIR /terabee/driver/post

COPY post/post.go /terabee/driver/post/