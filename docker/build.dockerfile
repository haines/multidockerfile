# syntax=docker/dockerfile:1.2

ARG \
  GO_VERSION


FROM \
  golang:${GO_VERSION}-alpine \
  AS go

ENV \
  CGO_ENABLED=0

WORKDIR \
  /multidockerfile


FROM \
  go \
  AS build

RUN \
  go install github.com/mitchellh/gox@latest


FROM \
  go \
  AS test
