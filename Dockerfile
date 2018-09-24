FROM golang:latest

COPY . /go/src/github.com/orenkay/matcha
WORKDIR /go/src/github.com/orenkay/matcha

RUN go get -t -v ./...
