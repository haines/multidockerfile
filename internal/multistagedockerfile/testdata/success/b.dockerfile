# syntax=docker/dockerfile:1.2

ARG \
  B

FROM \
  a \
  AS b

COPY \
  --from=c \
  ./ \
  ./
