#!/usr/bin/env bash
set -exu

export GO111MODULE=on
# installing golangci-lint as recommended on the project page
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
go mod download
golangci-lint run --disable typecheck --enable deadcode --enable varcheck --enable staticcheck
