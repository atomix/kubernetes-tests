FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/partition-group-membership/_output/bin/atomix-partition-group-member /usr/local/bin/atomix-partition-group-member

ENTRYPOINT ["atomix-partition-group-member"]
