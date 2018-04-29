FROM golang:1.9 as builder

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/jrnt30/k8-kms-enc-provider

COPY . .

RUN dep ensure && \
    CGO_ENABLED=0 go build --ldflags '-extldflags "-static"' -o k8-kms-enc-provider .

FROM alpine:3.7
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/jrnt30/k8-kms-enc-provider/k8-kms-enc-provider /usr/local/bin/k8-kms-enc-provider

ENTRYPOINT ["/usr/local/bin/k8-kms-enc-provider"]
