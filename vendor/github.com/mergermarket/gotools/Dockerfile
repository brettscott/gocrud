FROM golang:1.7

RUN go get -u github.com/alecthomas/gometalinter
RUN go get -u github.com/kardianos/govendor
RUN go get -u github.com/HewlettPackard/gas
RUN go get -u github.com/stretchr/testify/assert
RUN go get -u github.com/kyoh86/richgo
RUN gometalinter --install
WORKDIR /go/src/github.com/mergermarket/gotools
ADD . /go/src/github.com/mergermarket/gotools
CMD ./build.sh
