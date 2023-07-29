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

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    netcat-traditional vim \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* 

RUN curl -sSf https://temporal.download/cli.sh | sh

RUN rm -rf iwf-server && \
    make bins && \
    chmod +x /iwf/script/start-server.sh && \
    chmod +x /iwf/lite/start-lite-server.sh

EXPOSE  8801 7933 8233
CMD ["/iwf/lite/start-lite-server.sh"]
