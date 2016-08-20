#!/bin/sh

# Copyright (c) 2016 Paul Jolly <paul@myitcv.org.uk>, all rights reserved.
# Use of this document is governed by a license found in the LICENSE document.

set -e

go generate
go vet ./...

# no tests to run here...

cd cmd/sortGen/_testFiles/

go generate
go test
go vet
