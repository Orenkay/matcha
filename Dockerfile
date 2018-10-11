FROM golang:latest

ENV MATCHA_PATH  /go/src/github.com/orenkay/matcha

COPY . ${MATCHA_PATH}
WORKDIR ${MATCHA_PATH}

RUN go get -u github.com/onsi/ginkgo/ginkgo
RUN go get -u github.com/onsi/gomega/...

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
