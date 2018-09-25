FROM golang:latest

COPY . /go/src/github.com/orenkay/matcha
WORKDIR /go/src/github.com/orenkay/matcha

RUN apt-get update
RUN apt-get install -y sendmail

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
