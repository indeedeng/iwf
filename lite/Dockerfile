FROM alpine:latest AS deps

WORKDIR  /tmp/iwf-lite-deps

#get temporal cli for each architecture we support
ADD --checksum=sha256:e2f548af84e820b7d71f25fc6461a4e6a1ab7cdb6a80fde3c375641b0772375b --unpack=true https://github.com/temporalio/cli/releases/download/v1.7.2/temporal_cli_1.7.2_linux_amd64.tar.gz temporal/amd64/


ADD --checksum=sha256:0d60307cf036f5e29cf0d298dabc9dd62f1cce80f72cbfdbdb0e54ad7790a701 --unpack=true https://github.com/temporalio/cli/releases/download/v1.7.2/temporal_cli_1.7.2_linux_arm64.tar.gz temporal/arm64/

# Get minio source code
# last minio release we may want to switch to https://garagehq.deuxfleurs.fr/ or https://github.com/seaweedfs/seaweedfs
ADD --checksum=sha256:be6d0bd3696c3a13a35f02d3a0280b64319c67918b4501c5c3d87f96d000085c --unpack=false https://github.com/minio/minio/archive/refs/tags/RELEASE.2025-10-15T17-29-55Z.tar.gz minio/archive/minio.tar.gz

FROM golang:latest
ARG TARGETARCH
COPY . /iwf/
WORKDIR /iwf/

RUN apt-get update \
    && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    netcat-traditional vim \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* 

COPY --from=deps /tmp/iwf-lite-deps/temporal/${TARGETARCH}/temporal /usr/local/bin/temporal


# Build and Install MinIO
COPY --from=deps /tmp/iwf-lite-deps/minio/archive/minio.tar.gz /tmp/minio/archive/minio.tar.gz

RUN mkdir -p /tmp/minio/sourcecode/minio/

RUN tar -xf /tmp/minio/archive/minio.tar.gz --directory /tmp/minio/sourcecode/minio/

WORKDIR /tmp/minio/sourcecode/minio/minio-RELEASE.2025-10-15T17-29-55Z/
RUN go build /tmp/minio/sourcecode/minio/minio-RELEASE.2025-10-15T17-29-55Z/

RUN cp /tmp/minio/sourcecode/minio/minio-RELEASE.2025-10-15T17-29-55Z/minio /usr/local/bin/minio

WORKDIR /iwf/

RUN rm -rf iwf-server && \
    make bins && \
    chmod +x /iwf/lite/start-lite-server.sh

EXPOSE  8801 7233 8233 9000 9001
CMD ["/iwf/lite/start-lite-server.sh"]
