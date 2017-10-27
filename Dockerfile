FROM webhippie/alpine:latest
MAINTAINER Thomas Boerger <thomas@webhippie.de>

EXPOSE 8080 80 443
VOLUME ["/var/lib/oauth2-proxy"]

LABEL org.label-schema.version=latest
LABEL org.label-schema.name="OAuth2 Proxy"
LABEL org.label-schema.vendor="Thomas Boerger"
LABEL org.label-schema.schema-version="1.0"

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
    oauth2-proxy && \
  mkdir -p \
    /usr/share/oauth2-proxy

ENV OAUTH2_PROXY_SERVER_STORAGE /var/lib/oauth2-proxy
ENV OAUTH2_PROXY_SERVER_TEMPLATES /usr/share/oauth2-proxy/templates
ENV OAUTH2_PROXY_SERVER_ASSETS /usr/share/oauth2-proxy/assets

COPY assets /usr/share/oauth2-proxy/
COPY templates /usr/share/oauth2-proxy/
COPY oauth2-proxy /usr/bin/

USER oauth2-proxy
ENTRYPOINT ["/usr/bin/oauth2-proxy"]
CMD ["server"]
