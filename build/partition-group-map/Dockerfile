FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/partition-group-map/_output/bin/atomix-partition-group-map-member /usr/local/bin/atomix-partition-group-map-member

ENTRYPOINT ["atomix-partition-group-map-member"]
