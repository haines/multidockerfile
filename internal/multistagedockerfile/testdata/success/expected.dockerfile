# escape = `
# syntax = docker/dockerfile:1.2
ARG   A
ARG   B
ARG   C
FROM   alpine   AS a
FROM   alpine   AS c
FROM   alpine:d
FROM   alpine:e
FROM   a   AS b
COPY   --from=c   ./   ./
