FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/membership/_output/bin/atomix-member /usr/local/bin/atomix-member

ENTRYPOINT ["atomix-member"]
