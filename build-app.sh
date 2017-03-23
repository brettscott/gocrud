#!/bin/bash -e

set -o errexit
set -o nounset
set -o pipefail

if [ ! $(command -v gometalinter) ]
then
	go get github.com/alecthomas/gometalinter
	gometalinter --install --vendor
fi

gometalinter \
    --vendor \
	--exclude='error return value not checked.*(Close|Log|Print).*\(errcheck\)$' \
	--exclude='Errors unhandled.,LOW,HIGH' \
	--exclude='.*_test\.go:.*error return value not checked.*\(errcheck\)$' \
	--exclude='duplicate of.*_test.go.*\(dupl\)$' \
	--exclude='vendor' \
	--disable=aligncheck \
	--disable=gotype \
	--disable=structcheck \
	--disable=varcheck \
	--disable=unconvert \
	--disable=aligncheck \
	--disable=dupl \
	--disable=goconst  \
	--disable=gosimple  \
	--disable=staticcheck \
	--cyclo-over=30 \
	--tests \
	--deadline=150s \
	./...

go fmt $(go list ./... | grep -v /vendor/)
go test $(go list ./... | grep -v acceptance-tests | grep -v /vendor/) --cover -timeout 30s
go install