FROM golang:1.11

LABEL "name"="Go Runner"
LABEL "version"="0.1.0"

RUN \
  apt-get update && \
  apt-get install -y ca-certificates openssl zip && \
  update-ca-certificates && \
  rm -rf /var/lib/apt

ENV GO111MODULE on
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]