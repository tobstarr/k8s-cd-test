FROM golang:1.7.1-alpine

RUN apk add ca-certificates

RUN mkdir -p /go/src/github.com/tobstarr/k8s-cd-test
COPY . /go/src/github.com/tobstarr/k8s-cd-test

RUN cd /go/src/github.com/tobstarr/k8s-cd-test && \
    go build -o /usr/bin/k8s-cd-test

ENTRYPOINT ["/usr/bin/k8s-cd-test"]
