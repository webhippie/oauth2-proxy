FROM webhippie/alpine:latest

LABEL maintainer="Thomas Boerger <thomas@webhippie.de>" \
  org.label-schema.name="OAuth2 Proxy" \
  org.label-schema.vendor="Thomas Boerger" \
  org.label-schema.schema-version="1.0"

EXPOSE 80 443 9000
VOLUME ["/var/lib/oauth2-proxy"]

ENV OAUTH2_PROXY_HEALTH_ADDR 0.0.0.0:9000
ENV OAUTH2_PROXY_SERVER_STORAGE /var/lib/oauth2-proxy
ENV OAUTH2_PROXY_SERVER_TEMPLATES /usr/share/oauth2-proxy/templates
ENV OAUTH2_PROXY_SERVER_ASSETS /usr/share/oauth2-proxy/assets

ENTRYPOINT ["/usr/bin/oauth2-proxy"]
CMD ["server"]

RUN apk add --no-cache ca-certificates mailcap bash

COPY assets /usr/share/oauth2-proxy/
COPY templates /usr/share/oauth2-proxy/
COPY dist/binaries/oauth2-proxy-*-linux-amd64 /usr/bin/oauth2-proxy
