# syntax=docker/dockerfile:1.2

ARG \
  GO_VERSION


FROM \
  golang:${GO_VERSION}-alpine \
  AS build

WORKDIR \
  /build

COPY \
  go.mod \
  go.sum \
  ./

RUN \
  --mount="type=cache,id=go-mod-download,target=/go/pkg/mod/cache" \
  go mod download

COPY \
  ./ \
  ./

ARG \
  LDFLAGS

ENV \
  CGO_ENABLED=0

RUN \
  --mount="type=cache,id=go-build,target=/root/.cache/go-build" \
  go build \
    -ldflags "${LDFLAGS}" \
    -o target/multidockerfile \
    ./cmd/multidockerfile


FROM \
  scratch \
  AS multidockerfile

COPY \
  --from=build \
  /build/target/multidockerfile \
  /multidockerfile

ENTRYPOINT \
  ["/multidockerfile"]

CMD \
  ["--help"]
