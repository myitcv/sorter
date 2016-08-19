#!/bin/sh

set -e

go test ./...

cd cmd/sortGen/_testFiles/

go generate
go test
