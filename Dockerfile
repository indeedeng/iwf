# Build iwf-server binaries
FROM golang:1.17-alpine3.13 AS builder

ARG RELEASE_VERSION

RUN apk add --update --no-cache ca-certificates make git curl mercurial unzip
RUN apk add build-base

WORKDIR /iwf-server

# Making sure that dependency is not touched
ENV GOFLAGS="-mod=readonly"

# Copy go mod dependencies and build cache
COPY go.* ./
RUN go mod download

COPY . .
RUN rm -fr .bin .build

ENV CADENCE_NOTIFICATION_RELEASE_VERSION=$RELEASE_VERSION

RUN CGO_ENABLED=0 make bins

# Download dockerize
FROM alpine:3.11 AS dockerize

RUN apk add --no-cache openssl

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-alpine-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && echo "**** fix for host id mapping error ****" \
    && chown root:root /usr/local/bin/dockerize

RUN apk add --update --no-cache ca-certificates tzdata bash curl
RUN test ! -e /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf
SHELL [ "/bin/bash", "-c" ]

# Cadence server
FROM alpine AS iwf-server

ENV CADENCE_NOTIFICATION_HOME=/etc/iwf-server
RUN mkdir -p /etc/iwf-server

COPY --from=builder /iwf-server/iwf-server /usr/local/bin
COPY --from=builder /iwf-server/config /iwf-server/config
COPY --from=dockerize /usr/local/bin/dockerize /usr/local/bin

COPY /config/config_template.yaml /etc/iwf-server/config
COPY /start.sh /start.sh

WORKDIR /etc/iwf-server

ENV SERVICES="api,interpreter"
RUN chmod +x /start.sh
CMD /start.sh