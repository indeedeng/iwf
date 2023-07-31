FROM golang:latest

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
    chmod +x /iwf/lite/start-lite-server.sh

EXPOSE  8801 7233 8233
CMD ["/iwf/lite/start-lite-server.sh"]
