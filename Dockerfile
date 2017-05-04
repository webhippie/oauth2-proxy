FROM webhippie/alpine:latest
MAINTAINER Thomas Boerger <thomas@webhippie.de>

EXPOSE 8080
VOLUME ["/var/lib/oauth2-proxy"]

RUN apk update && \
  apk add \
    ca-certificates \
    bash && \
  rm -rf \
    /var/cache/apk/* && \
  addgroup \
    -g 1000 \
    oauth2-proxy && \
  adduser -D \
    -h /var/lib/oauth2-proxy \
    -s /bin/bash \
    -G oauth2-proxy \
    -u 1000 \
    oauth2-proxy

COPY oauth2-proxy /usr/bin/

USER oauth2-proxy
ENTRYPOINT ["/usr/bin/oauth2-proxy"]
CMD ["server"]

# ARG VERSION
# ARG BUILD_DATE
# ARG VCS_REF

# LABEL org.label-schema.version=$VERSION
# LABEL org.label-schema.build-date=$BUILD_DATE
# LABEL org.label-schema.vcs-ref=$VCS_REF
LABEL org.label-schema.vcs-url="https://github.com/webhippie/oauth2-proxy.git"
LABEL org.label-schema.name="OAuth2 Proxy"
LABEL org.label-schema.vendor="Thomas Boerger"
LABEL org.label-schema.schema-version="1.0"
