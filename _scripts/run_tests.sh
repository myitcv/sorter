#!/bin/sh

set -e

go generate

# no tests to run here...

cd cmd/sortGen/_testFiles/

go generate
go test
