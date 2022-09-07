FROM golang:1.10 AS build
WORKDIR /go/src
COPY workflow ./workflow
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o workflow .

FROM scratch AS runtime
ENV GIN_MODE=release
COPY --from=build /go/src/workflow ./
EXPOSE 8080/tcp
ENTRYPOINT ["./workflow"]
