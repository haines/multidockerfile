FROM \
  alpine \
  AS a

COPY \
  --from=c \
  /foo \
  ./


FROM \
  a \
  AS b

COPY \
  --from=c \
  /bar \
  ./


FROM \
  b \
  AS c
