FROM golang:1.8

RUN apt-get update
RUN apt-get install libsasl2-dev

RUN go get -u github.com/alecthomas/gometalinter
RUN go get -u github.com/kardianos/govendor
RUN gometalinter --install
WORKDIR /go/src/github.com/brettscott/gocrud
ADD . /go/src/github.com/brettscott/gocrud
RUN go build
EXPOSE 8080
