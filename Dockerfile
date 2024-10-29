FROM alpine:3.20

RUN apk add --no-cache tini

# Include curl in the final image.
RUN set -ex \
    && apk update \
    && apk add --no-cache curl tini \
    && rm -rf /var/cache/apk/*  \
    && rm -rf /tmp/*

COPY smd /
COPY smd-loader /
COPY smd-init /
RUN mkdir /persistent_migrations
COPY migrations/* /persistent_migrations/

EXPOSE 27779

# nobody 65534:65534
USER 65534:65534

CMD [ "/smd" ]

ENTRYPOINT [ "/sbin/tini", "--" ]
