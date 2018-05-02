#!/usr/bin/env bash

# Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
# Use of this document is governed by a license found in the LICENSE document.

source "${BASH_SOURCE%/*}/common.bash"

export GOPATH=$PWD/_vendor:$GOPATH
export PATH="${GOPATH//://bin:}/bin:$PATH"

rm -f !(_vendor)/**/gen_*.go

# TODO use gg
go install myitcv.io/sorter/cmd/sortGen
go install myitcv.io/immutable/cmd/immutableGen

go generate ./...
go install ./...
go vet ./...
go test ./...

cd cmd/sortGen/_testFiles/

go generate ./...
go test ./...
go install ./...
go vet ./...
