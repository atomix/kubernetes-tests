FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add libc6-compat

USER nobody

ADD build/peer/_output/bin/atomix-peer /usr/local/bin/atomix-peer

ENTRYPOINT ["atomix-peer"]