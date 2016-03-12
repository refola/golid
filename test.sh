#!/bin/sh

echo "Processing code."
# gofmt everything
find . -type f -name \*.go -print0 | xargs -0 -r -n1 gofmt -w
# reinstall updated binary
go install github.com/refola/piklisp_go/piklisp

echo "Running tests."
go test ./parse/
