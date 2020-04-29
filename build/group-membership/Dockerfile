FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/group-membership/_output/bin/atomix-group-member /usr/local/bin/atomix-group-member

ENTRYPOINT ["atomix-group-member"]
