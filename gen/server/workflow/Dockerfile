FROM golang:1.10 AS build
WORKDIR /go/src
COPY serverapi ./serverapi
COPY main.go .

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN go build -a -installsuffix cgo -o serverapi .

FROM scratch AS runtime
COPY --from=build /go/src/serverapi ./
EXPOSE 8080/tcp
ENTRYPOINT ["./serverapi"]
