FROM golang:alpine as build-env

RUN apk add build-base linux-headers
RUN apk --update add ca-certificates

RUN mkdir /external-initiator
WORKDIR /external-initiator
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.15.1/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 379c61d2e53f87f63639eaa8ba8bbe687e5833158bb10e0dc0523c3377248d01

# Delete ./integration folder that is not needed in the context of external-initiator,
# but is required in the context of mock-client build.
RUN rm -rf ./integration
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags "muslc make build" -a -installsuffix cgo -o /go/bin/external-initiator

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/bin/external-initiator /go/bin/external-initiator

EXPOSE 8080

ENTRYPOINT ["/go/bin/external-initiator"]
