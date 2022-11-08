FROM golang:latest

ENV CONFIG_TEMPLATE_PATH=''
ENV SRC_ROOT=''
ENV BACKEND_DEPENDENCY=${BACKEND_DEPENDENCY:-'temporal'}
ENV CONFIG_TEMPLATE_PATH=''
ENV SRC_ROOT=''
ENV HOST=''
ENV TEMPORAL_SERVICE_NAME=''
ENV CADENCE_SERVICE_NAME=''

COPY . /iwf/
WORKDIR /iwf/

RUN rm -rf iwf-server  \
    && apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    netcat \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* && \
    make bins && \
    chmod +x /iwf/script/start-server.sh

EXPOSE  8801
ENTRYPOINT ["/iwf/script/start-server.sh"]
