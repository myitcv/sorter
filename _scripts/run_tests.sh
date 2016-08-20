#!/bin/sh

set -e

go generate
go vet ./...

# no tests to run here...

cd cmd/sortGen/_testFiles/

go generate
go test
go vet
